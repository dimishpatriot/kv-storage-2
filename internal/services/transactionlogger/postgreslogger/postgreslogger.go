package postgreslogger

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"

	"github.com/dimishpatriot/kv-storage/internal/services/transactionlogger"
	"github.com/dimishpatriot/kv-storage/internal/storage/postgresstorage"
)

type PostgresTransactionLogger struct {
	events  chan<- transactionlogger.Event
	errors  <-chan error
	db      *sql.DB
	logger  *log.Logger
	storage *postgresstorage.PostgresStorage
}

type PostgresDBParams struct {
	DBName   string
	Host     string
	User     string
	Password string
	SSLMode  string
}

func New(
	logger *log.Logger,
	dbParams PostgresDBParams,
) (transactionlogger.TransactionLogger, *sql.DB, error) {
	connString := fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=%s",
		dbParams.User, dbParams.Password, dbParams.Host, dbParams.DBName, dbParams.SSLMode,
	)
	db, err := getDBConnection(connString)
	if err != nil {
		return nil, nil, fmt.Errorf("cant get db: %w", err)
	}
	storage := postgresstorage.New(db, "transactions")

	if exists := storage.VerifyTableExists(); !exists {
		err = storage.CreateTable()
		if err != nil {
			return nil, nil, fmt.Errorf("can't create table: %w", err)
		}
	}

	return &PostgresTransactionLogger{
		logger:  logger,
		db:      db,
		storage: storage,
	}, db, nil
}

func getDBConnection(connStr string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open db: %w", err)
	}
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to open db connection: %w", err)
	}
	return db, nil
}

func (l *PostgresTransactionLogger) Run() {
	var err error

	l.logger.Println("dataLogger run...")

	events := make(chan transactionlogger.Event, 16)
	l.events = events
	errors := make(chan error, 1)
	l.errors = errors

	go func() {
		for event := range events {
			switch event.EventType {
			case transactionlogger.EventPut:
				if err = l.storage.Put(event.Key, event.Value); err != nil {
					errors <- err
					return
				}
			case transactionlogger.EventDelete:
				if err = l.storage.Delete(event.Key); err != nil {
					errors <- err
					return
				}
			}
		}
	}()
}

func (l *PostgresTransactionLogger) ReadEvents() (<-chan transactionlogger.Event, <-chan error) {
	l.logger.Println("read events...")

	outEvent := make(chan transactionlogger.Event)
	outError := make(chan error, 1)

	go func() {
		defer close(outEvent)
		defer close(outError)

		res, err := l.storage.GetAll()
		if err != nil {
			outError <- err
		}
		for _, r := range res {
			outEvent <- r
		}
	}()

	return outEvent, outError
}

func (l *PostgresTransactionLogger) WritePut(key, value string) {
	l.logger.Printf("write put: {%s: %s}", key, value)

	l.events <- transactionlogger.Event{
		EventType: transactionlogger.EventPut, Key: key, Value: value,
	}
}

func (l *PostgresTransactionLogger) WriteDelete(key string) {
	l.logger.Printf("write delete {%s}", key)

	l.events <- transactionlogger.Event{
		EventType: transactionlogger.EventDelete, Key: key,
	}
}

func (l *PostgresTransactionLogger) Err() <-chan error {
	l.logger.Println("getting errors channel...")

	return l.errors
}
