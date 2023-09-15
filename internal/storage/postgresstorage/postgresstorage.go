package postgresstorage

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/dimishpatriot/kv-storage/internal/services/transactionlogger"
	"github.com/dimishpatriot/kv-storage/internal/storage"
)

type PostgresStorage struct {
	db *sql.DB
}

func New(db *sql.DB) *PostgresStorage {
	return &PostgresStorage{db: db}
}

func (s *PostgresStorage) VerifyTableExists() bool {
	qs := "SELECT * FROM transactions LIMIT 1;"
	if _, err := s.db.Query(qs); err != nil {
		return false
	}

	return true
}

func (s *PostgresStorage) CreateTable() error {
	q := `
	CREATE TABLE transactions (
	sequence BIGSERIAL PRIMARY KEY,
	event_type SMALLINT,
	key TEXT,
	value TEXT)`
	if _, err := s.db.Exec(q); err != nil {
		return fmt.Errorf("can't create transactions table: %w", err)
	}

	return nil
}

func (s *PostgresStorage) Put(k, v string) error {
	q := `
	INSERT INTO transactions 
	(event_type, key, value) 
	VALUES ($1, $2, $3)
`
	if _, err := s.db.Exec(q, transactionlogger.EventPut, k, v); err != nil {
		return fmt.Errorf("failed to insert data: %w", err)
	}

	return nil
}

type Row struct {
	Sequence  string
	EventType transactionlogger.EventType
	Key       string
	Value     string
}

func (s *PostgresStorage) GetAll() ([]transactionlogger.Event, error) {
	q := "SELECT * FROM transactions ORDER BY sequence"
	result := []transactionlogger.Event{}

	rows, err := s.db.Query(q)
	if err != nil {
		return nil, fmt.Errorf("get all events error: %w", err)
	}
	defer rows.Close()

	e := transactionlogger.Event{}
	for rows.Next() {
		err = rows.Scan(&e.Sequence, &e.EventType, &e.Key, &e.Value)
		if err != nil {
			return nil, fmt.Errorf("error reading row: %w", err)
		}
		result = append(result, e)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("fail to read transaction log: %w", err)
	}

	return result, nil
}

func (s *PostgresStorage) Get(k string) (string, error) {
	q := `
	SELECT event_type, key, value 
	FROM transactions 
	WHERE key=$1`
	row := s.db.QueryRow(q, k)

	e := transactionlogger.Event{}
	err := row.Scan(&e.EventType, &e.Key, &e.Value)

	if errors.Is(err, sql.ErrNoRows) {
		return "", storage.ErrorNoSuchKey
	}

	return e.Value, nil
}

func (s *PostgresStorage) Delete(k string) error {
	q := `
	DELETE FROM transactions 
	WHERE key=$1`

	_, err := s.db.Exec(q, k)
	if err != nil {
		return fmt.Errorf("failed to clear data: %w", err)
	}

	return nil
}
