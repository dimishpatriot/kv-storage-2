package dlhandler_test

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/dimishpatriot/kv-storage/internal/handler"
	"github.com/dimishpatriot/kv-storage/internal/handler/dlhandler"
	"github.com/dimishpatriot/kv-storage/internal/storage"
	"github.com/dimishpatriot/kv-storage/internal/storage/localstorage"
	"github.com/dimishpatriot/kv-storage/internal/transactionlogger"
	"github.com/dimishpatriot/kv-storage/internal/transactionlogger/mocks"
	"github.com/gorilla/mux"
)

var (
	logger     *log.Logger
	dataLogger transactionlogger.TransactionLogger
	store      storage.Storage
	dlh        handler.Handler
)

func TestMain(m *testing.M) {
	logger = log.New(os.Stdout, "INFO:", log.Lshortfile|log.Ltime|log.Lmicroseconds|log.Ldate)
	dataLogger = &mocks.MockDataLogger{}

	os.Exit(m.Run())
}

func setupTest(tb testing.TB) func(tb testing.TB) {
	store = localstorage.New()
	dlh = dlhandler.New(logger, dataLogger, store)

	// prepare test data
	_ = store.Put("1", "ONE")
	_ = store.Put("-1", "minus ONE")
	_ = store.Put("!@#$*()_+><", "symbols")
	_ = store.Put("0123456789", "numbers")

	return func(tb testing.TB) {
		// run after each test
		// ...
	}
}

func getPath(key string) string {
	return fmt.Sprintf("/v1/%s", key)
}

func TestDataLoggerHandler_Put(t *testing.T) {
	type args struct {
		key   string
		value string
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
	}{
		{name: "success put by new key", args: args{key: "new", value: "new value"}, wantStatus: http.StatusCreated},
		{name: "success put by existing key", args: args{key: "1", value: "one new"}, wantStatus: http.StatusCreated},
		{name: "failed put by empty key", args: args{key: "", value: "one new"}, wantStatus: http.StatusBadRequest},
		{name: "failed put by empty value", args: args{key: "key", value: ""}, wantStatus: http.StatusBadRequest},
		{name: "failed put by long key", args: args{key: "12345678901234567890123456789012345678901234567890123456789012345", value: ""}, wantStatus: http.StatusBadRequest},
		{name: "failed put by long value", args: args{key: "key", value: "1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890"}, wantStatus: http.StatusBadRequest},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			after := setupTest(t)
			defer after(t)

			res := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPut, getPath(tt.args.key), strings.NewReader(tt.args.value))
			// setup url vars for mux.Vars()
			r = mux.SetURLVars(r,
				map[string]string{
					"key": tt.args.key,
				})

			dlh.Put(res, r)
			if res.Code != tt.wantStatus {
				t.Errorf("got status %d, wont %d", res.Code, tt.wantStatus)
			}
		})
	}
}

func TestDataLoggerHandler_Get(t *testing.T) {
	type args struct {
		key string
	}
	type want struct {
		status int
		value  string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{name: "success get by existing key", args: args{key: "1"}, want: want{value: "ONE", status: http.StatusOK}},
		{name: "success get by existing symbol+key", args: args{key: "-1"}, want: want{value: "minus ONE", status: http.StatusOK}},
		{name: "success get by symbolic key", args: args{key: "!@#$*()_+><"}, want: want{value: "symbols", status: http.StatusOK}},
		{name: "failed get by existing key", args: args{key: "11"}, want: want{value: "", status: http.StatusNotFound}},
		{name: "failed get by empty key", args: args{key: ""}, want: want{value: "", status: http.StatusBadRequest}},
		{name: "failed get by long key", args: args{key: "12345678901234567890123456789012345678901234567890123456789012345"}, want: want{value: "", status: http.StatusBadRequest}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			after := setupTest(t)
			defer after(t)

			res := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, getPath(tt.args.key), nil)
			// setup url vars for mux.Vars()
			r = mux.SetURLVars(r,
				map[string]string{
					"key": tt.args.key,
				})

			dlh.Get(res, r)

			if res.Code != tt.want.status {
				t.Errorf("got status %d, wont %d", res.Code, tt.want.status)
			}
			if res.Code == http.StatusOK {
				if res.Body.String() != tt.want.value {
					t.Errorf("value got=%s, want=%s", res.Body.String(), tt.want.value)
				}
			}
		})
	}
}

func TestDataLoggerHandler_Delete(t *testing.T) {
	type args struct {
		key string
	}
	type want struct {
		status int
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{name: "success delete by existing key", args: args{key: "1"}, want: want{status: http.StatusOK}},
		{name: "success delete by existing symbol+key", args: args{key: "-1"}, want: want{status: http.StatusOK}},
		{name: "success delete by symbolic key", args: args{key: "!@#$*()_+><"}, want: want{status: http.StatusOK}},
		{name: "failed delete by existing key", args: args{key: "11"}, want: want{status: http.StatusNotFound}},
		{name: "failed delete by empty key", args: args{key: ""}, want: want{status: http.StatusBadRequest}},
		{name: "failed delete by long key", args: args{key: "12345678901234567890123456789012345678901234567890123456789012345"}, want: want{status: http.StatusBadRequest}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			after := setupTest(t)
			defer after(t)

			res := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodDelete, getPath(tt.args.key), nil)
			// setup url vars for mux.Vars()
			r = mux.SetURLVars(r,
				map[string]string{
					"key": tt.args.key,
				})

			dlh.Delete(res, r)

			if res.Code != tt.want.status {
				t.Errorf("got status %d, wont %d", res.Code, tt.want.status)
			}
		})
	}
}
