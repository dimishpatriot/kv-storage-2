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

	// need add sequence for sqlite3 test base!
	q = fmt.Sprintf(`
	INSERT INTO %s 
	(sequence, event_type, key, value) 
	VALUES ($1, $2, $3, $4)
	`, tableName)
	_, _ = db.Exec(q, 1, transactionlogger.EventPut, "one", "ONE")
	_, _ = db.Exec(q, 2, transactionlogger.EventPut, "2", "two")

	return m.Run(), nil
}

func TestPostgresStorage_VerifyTableExists(t *testing.T) {
	type test struct {
		name      string
		tableName string
		isExists  bool
	}
	tests := []test{
		{
			"existing table",
			tableName,
			true,
		},
		{
			"not existing table",
			"ABC",
			false,
		},
		{
			"empty table name",
			"",
			false,
		},
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
		{
			"new table",
			"new",
			false,
		},
		{
			"existing table",
			tableName,
			true,
		},
		{
			"table with empty name",
			"",
			true,
		},
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

func TestPostgresStorage_GetAll(t *testing.T) {
	type want struct {
		result []transactionlogger.Event
		err    bool
	}
	type test struct {
		name string
		want want
	}
	tests := []test{
		{
			"simple",
			want{
				[]transactionlogger.Event{
					{Sequence: 1, EventType: transactionlogger.EventPut, Key: "one", Value: "ONE"},
					{Sequence: 2, EventType: transactionlogger.EventPut, Key: "2", Value: "two"},
				},
				false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := postgresstorage.New(db, tableName)
			res, err := store.GetAll()

			assert.Equal(t, tt.want.err, err != nil)
			assert.EqualValues(t, tt.want.result, res)
		})
	}
}

func TestPostgresStorage_Put(t *testing.T) {
	type args struct {
		key   string
		value string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"correct args",
			args{key: "key", value: "value"},
			false,
		},
		{
			"equal args",
			args{key: "equal", value: "equal"},
			false,
		},
		{
			"short args",
			args{key: "k", value: "v"},
			false,
		},
		{
			"with numbers args",
			args{key: "100", value: "500"},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := postgresstorage.New(db, tableName)

			err := s.Put(tt.args.key, tt.args.value)

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
		{
			"existing key",
			"one",
			want{value: "ONE", error: false},
		},
		{
			"no existing key",
			" one ",
			want{value: "", error: true},
		},
		{
			"empty key",
			"",
			want{value: "", error: true},
		},
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

func TestPostgresStorage_Delete(t *testing.T) {
	type want struct {
		error bool
	}
	tests := []struct {
		name string
		key  string
		want want
	}{
		{
			"existing key",
			"one",
			want{error: false},
		},
		{
			"no existing key",
			" one ",
			want{error: true},
		},
		{
			"empty key",
			"",
			want{error: true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := postgresstorage.New(db, tableName)

			err := s.Delete(tt.key)
			assert.Equal(t, tt.want.error, err != nil)
		})
	}
}
