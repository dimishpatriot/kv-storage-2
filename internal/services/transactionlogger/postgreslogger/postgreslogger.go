package postgreslogger

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"

	"github.com/dimishpatriot/kv-storage/internal/services/transactionlogger"
	"github.com/dimishpatriot/kv-storage/internal/storage"
)

type PostgresTransactionLogger struct {
	events  chan<- transactionlogger.Event
	errors  <-chan error
	db      *sql.DB
	logger  *log.Logger
	storage storage.Storage
}

type PostgresDBParams struct {
	name     string
	host     string
	user     string
	password string
}

func New(
	logger *log.Logger,
	dbParams PostgresDBParams,
	storage storage.Storage,
) (transactionlogger.TransactionLogger, error) {
	connStr := fmt.Sprintf(
		"host=%s dbname=%s user=%s password=%s",
		dbParams.host, dbParams.name, dbParams.user, dbParams.password,
	)

	db, err := getDBConnection(connStr)
	if err != nil {
		return nil, fmt.Errorf("cant get db: %w", err)
	}

	ptl := PostgresTransactionLogger{logger: logger, db: db, storage: storage}

	exists, err := ptl.verifyTableExists()
	if err != nil {
		return nil, fmt.Errorf("cant verify table exists: %w", err)
	}
	if !exists {
		if err = ptl.createTable(); err != nil {
			return nil, fmt.Errorf("cant create table: %w", err)
		}
	}

	return &ptl, nil
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

func (l *PostgresTransactionLogger) verifyTableExists() (bool, error) {
	// TODO:

	return true, nil
}

func (l *PostgresTransactionLogger) createTable() error {
	// TODO:

	return nil
}

func (l *PostgresTransactionLogger) Run() {
	l.logger.Println("dataLogger run...")

	events := make(chan transactionlogger.Event, 16)
	l.events = events
	errors := make(chan error, 1)
	l.errors = errors

	go func() {
		// TODO:
	}()
}

func (l *PostgresTransactionLogger) ReadEvents() (<-chan transactionlogger.Event, <-chan error) {
	l.logger.Println("read events...")

	outEvent := make(chan transactionlogger.Event)
	outError := make(chan error, 1)

	go func() {
		defer close(outEvent)
		defer close(outError)

		// TODO:
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
