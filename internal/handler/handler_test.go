package handler_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dimishpatriot/kv-storage/internal/handler"
	"github.com/dimishpatriot/kv-storage/internal/services/keyservice"
	"github.com/dimishpatriot/kv-storage/internal/storage"
	"github.com/gorilla/mux"
)

var (
	dlh         handler.Handler
	serviceMock *keyservice.MockKeyService
)

func setupTest(tb testing.TB) func(tb testing.TB) {
	serviceMock = keyservice.NewMockKeyService(tb)
	dlh = handler.New(serviceMock)

	return func(tb testing.TB) {
		// run after each test
		// ...
	}
}

func getPath(key string) string {
	return fmt.Sprintf("/v1/%s", key)
}

func TestDataHandler_Put(t *testing.T) {
	type args struct {
		key   string
		value string
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
	}{
		{
			"success put by new key",
			args{key: "new", value: "new value"},
			http.StatusCreated,
		},
		{
			"success put by existing key",
			args{key: "1", value: "one new"},
			http.StatusCreated,
		},
		{
			"failed put by empty key",
			args{key: "", value: "one new"},
			http.StatusBadRequest,
		},
		{
			"failed put by empty value",
			args{key: "key", value: ""},
			http.StatusBadRequest,
		},
		{
			"failed put by long key",
			args{
				key:   "12345678901234567890123456789012345678901234567890123456789012345",
				value: "",
			},
			http.StatusBadRequest,
		},
		{
			"failed put by long value",
			args{
				key:   "key",
				value: "1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890",
			},
			http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			after := setupTest(t)
			defer after(t)

			if tt.wantStatus == http.StatusCreated {
				serviceMock.EXPECT().Put(tt.args.key, tt.args.value).Return(nil)
			}

			res := httptest.NewRecorder()
			r := httptest.NewRequest(
				http.MethodPut,
				getPath(tt.args.key),
				strings.NewReader(tt.args.value),
			)
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

func TestDataHandler_Get(t *testing.T) {
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
		{
			"success get by existing key",
			args{key: "1"},
			want{value: "ONE", status: http.StatusOK},
		},
		{
			"success get by existing symbol+key",
			args{key: "-1"},
			want{value: "minus ONE", status: http.StatusOK},
		},
		{
			"success get by symbolic key",
			args{key: "!@#$*()_+><"},
			want{value: "symbols", status: http.StatusOK},
		},
		{
			"failed get by existing key",
			args{key: "11"},
			want{value: "", status: http.StatusNotFound},
		},
		{
			"failed get by empty key",
			args{key: ""},
			want{value: "", status: http.StatusBadRequest},
		},
		{
			"failed get by long key",
			args{key: "12345678901234567890123456789012345678901234567890123456789012345"},
			want{value: "", status: http.StatusBadRequest},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			after := setupTest(t)
			defer after(t)

			if tt.want.status == http.StatusOK {
				serviceMock.EXPECT().Get(tt.args.key).Return(tt.want.value, nil)
			}
			if tt.want.status == http.StatusNotFound {
				serviceMock.EXPECT().Get(tt.args.key).Return("", storage.ErrorNoSuchKey)
			}

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

func TestDataHandler_Delete(t *testing.T) {
	type args struct {
		key string
	}

	tests := []struct {
		name       string
		args       args
		wantStatus int
	}{
		{
			"success delete by existing key",
			args{key: "1"},
			http.StatusOK,
		},
		{
			"success delete by existing symbol+key",
			args{key: "-1"},
			http.StatusOK,
		},
		{
			"success delete by symbolic key",
			args{key: "!@#$*()_+><"},
			http.StatusOK,
		},
		{
			"failed delete by existing key",
			args{key: "11"},
			http.StatusNotFound,
		},
		{
			"failed delete by empty key",
			args{key: ""},
			http.StatusBadRequest,
		},
		{
			"failed delete by long key",
			args{
				"12345678901234567890123456789012345678901234567890123456789012345",
			},
			http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			after := setupTest(t)
			defer after(t)
			if tt.wantStatus == http.StatusOK {
				serviceMock.EXPECT().Delete(tt.args.key).Return(nil)
			}
			if tt.wantStatus == http.StatusNotFound {
				serviceMock.EXPECT().Delete(tt.args.key).Return(storage.ErrorNoSuchKey)
			}

			res := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodDelete, getPath(tt.args.key), nil)
			// setup url vars for mux.Vars()
			r = mux.SetURLVars(r,
				map[string]string{
					"key": tt.args.key,
				})

			dlh.Delete(res, r)

			if res.Code != tt.wantStatus {
				t.Errorf("got status %d, wont %d", res.Code, tt.wantStatus)
			}
		})
	}
}
