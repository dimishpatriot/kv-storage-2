package postgresstorage

import (
	"database/sql"

	"github.com/dimishpatriot/kv-storage/internal/storage"
)

type PostgresStorage struct {
	db *sql.DB
}

func New(db *sql.DB) storage.Storage {
	return &PostgresStorage{db: db}
}

func (ls *PostgresStorage) Put(k string, v string) error {
	// TODO:

	return nil
}

func (ls *PostgresStorage) Get(k string) (string, error) {
	// TODO:

	return "", nil
}

func (ls *PostgresStorage) Delete(k string) error {
	// TODO:

	return nil
}
