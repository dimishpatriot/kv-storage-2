package dlhandler_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dimishpatriot/kv-storage/internal/handler"
	"github.com/dimishpatriot/kv-storage/internal/handler/dlhandler"
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
	dlh = dlhandler.New(serviceMock)

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
			name:       "success put by new key",
			args:       args{key: "new", value: "new value"},
			wantStatus: http.StatusCreated,
		},
		{
			name:       "success put by existing key",
			args:       args{key: "1", value: "one new"},
			wantStatus: http.StatusCreated,
		},
		{
			name:       "failed put by empty key",
			args:       args{key: "", value: "one new"},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "failed put by empty value",
			args:       args{key: "key", value: ""},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "failed put by long key",
			args: args{
				key:   "12345678901234567890123456789012345678901234567890123456789012345",
				value: "",
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "failed put by long value",
			args: args{
				key:   "key",
				value: "1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890",
			},
			wantStatus: http.StatusBadRequest,
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
			name: "success get by existing key",
			args: args{key: "1"},
			want: want{value: "ONE", status: http.StatusOK},
		},
		{
			name: "success get by existing symbol+key",
			args: args{key: "-1"},
			want: want{value: "minus ONE", status: http.StatusOK},
		},
		{
			name: "success get by symbolic key",
			args: args{key: "!@#$*()_+><"},
			want: want{value: "symbols", status: http.StatusOK},
		},
		{
			name: "failed get by existing key",
			args: args{key: "11"},
			want: want{value: "", status: http.StatusNotFound},
		},
		{
			name: "failed get by empty key",
			args: args{key: ""},
			want: want{value: "", status: http.StatusBadRequest},
		},
		{
			name: "failed get by long key",
			args: args{key: "12345678901234567890123456789012345678901234567890123456789012345"},
			want: want{value: "", status: http.StatusBadRequest},
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
			name:       "success delete by existing key",
			args:       args{key: "1"},
			wantStatus: http.StatusOK,
		},
		{
			name:       "success delete by existing symbol+key",
			args:       args{key: "-1"},
			wantStatus: http.StatusOK,
		},
		{
			name:       "success delete by symbolic key",
			args:       args{key: "!@#$*()_+><"},
			wantStatus: http.StatusOK,
		},
		{
			name:       "failed delete by existing key",
			args:       args{key: "11"},
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "failed delete by empty key",
			args:       args{key: ""},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "failed delete by long key",
			args: args{
				key: "12345678901234567890123456789012345678901234567890123456789012345",
			},
			wantStatus: http.StatusBadRequest,
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
