package postgresstorage_test

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/dimishpatriot/kv-storage/internal/services/transactionlogger"
	"github.com/dimishpatriot/kv-storage/internal/storage/postgresstorage"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

var (
	db        *sql.DB
	tableName = "transactions"
)

func TestMain(m *testing.M) {
	code, err := prepareDB(m)
	if err != nil {
		log.Fatalf("can't prepare test DB: %s", err)
	}

	os.Exit(code)
}

func prepareDB(m *testing.M) (code int, err error) {
	var q string

	db, err = sql.Open("sqlite3", "file:./test.db?cache=shared")
	if err != nil {
		return -1, fmt.Errorf("can't connect to db: %w", err)
	}

	defer func() {
		q = fmt.Sprintf("DELETE FROM %s", tableName)
		_, _ = db.Exec(q)
		db.Close()
	}()

	q = fmt.Sprintf(`
	CREATE TABLE %s (
	sequence BIGSERIAL PRIMARY KEY,
	event_type SMALLINT,
	key TEXT NOT NULL,
	value TEXT NOT NULL)
	`, tableName)
	_, _ = db.Exec(q)

	q = fmt.Sprintf(`
	INSERT INTO %s 
	(event_type, key, value) 
	VALUES ($1, $2, $3)
	`, tableName)
	_, _ = db.Exec(q, transactionlogger.EventPut, "one", "ONE")
	_, _ = db.Exec(q, transactionlogger.EventPut, "2", "two")

	return m.Run(), nil
}

func TestPostgresStorage_VerifyTableExists(t *testing.T) {
	type test struct {
		name      string
		tableName string
		isExists  bool
	}
	tests := []test{
		{name: "existing table", tableName: tableName, isExists: true},
		{name: "not existing table", tableName: "ABC", isExists: false},
		{name: "empty table name", tableName: "", isExists: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := postgresstorage.New(db, tt.tableName)

			exists := s.VerifyTableExists()

			assert.Equal(t, exists, tt.isExists)
		})
	}
}

func TestPostgresStorage_CreateTable(t *testing.T) {
	type test struct {
		name      string
		tableName string
		wantErr   bool
	}
	tests := []test{
		{name: "new table", tableName: "new", wantErr: false},
		{name: "existing table", tableName: tableName, wantErr: true},
		{name: "table with empty name", tableName: "", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := postgresstorage.New(db, tt.tableName)

			err := s.CreateTable()
			defer func() {
				if tt.tableName != tableName {
					q := fmt.Sprintf("DROP TABLE %s", tt.tableName)
					_, _ = db.Exec(q)
				}
			}()

			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestPostgresStorage_Put(t *testing.T) {
	type args struct {
		k string
		v string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "correct args", args: args{k: "key", v: "value"}, wantErr: false},
		{name: "equal args", args: args{k: "equal", v: "equal"}, wantErr: false},
		{name: "short args", args: args{k: "k", v: "v"}, wantErr: false},
		{name: "with numbers args", args: args{k: "100", v: "500"}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := postgresstorage.New(db, tableName)

			err := s.Put(tt.args.k, tt.args.v)

			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestPostgresStorage_Get(t *testing.T) {
	type want struct {
		value string
		error bool
	}
	tests := []struct {
		name string
		key  string
		want want
	}{
		{name: "existing key", key: "one", want: want{value: "ONE", error: false}},
		{name: "no existing key", key: " one ", want: want{value: "", error: true}},
		{name: "empty key", key: "", want: want{value: "", error: true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := postgresstorage.New(db, tableName)

			got, err := s.Get(tt.key)
			assert.Equal(t, tt.want.value, got)
			assert.Equal(t, tt.want.error, err != nil)
		})
	}
}
