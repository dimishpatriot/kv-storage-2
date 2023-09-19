package keyservice_test

import (
	"errors"
	"io"
	"log"
	"testing"

	"github.com/dimishpatriot/kv-storage/internal/services/keyservice"
	"github.com/dimishpatriot/kv-storage/internal/services/transactionlogger"
	"github.com/dimishpatriot/kv-storage/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	srv         keyservice.KeyService
	logger      *log.Logger
	storageMock *storage.MockStorage
	tLoggerMock *transactionlogger.MockTransactionLogger
)

func setupTest(tb testing.TB) func(tb testing.TB) {
	logger = log.New(io.Discard, "", log.Lshortfile|log.Ltime|log.Lmicroseconds|log.Ldate)
	storageMock = storage.NewMockStorage(tb)
	tLoggerMock = transactionlogger.NewMockTransactionLogger(tb)
	srv = keyservice.New(logger, storageMock, tLoggerMock)

	return func(tb testing.TB) {
	}
}

func TestKeyService_Put(t *testing.T) {
	type args struct {
		key   string
		value string
	}
	type want struct {
		err bool
	}
	type test struct {
		name string
		args args
		want want
	}
	tests := []test{
		{
			"simple args",
			args{key: "one", value: "1"},
			want{err: false},
		},
		{
			"long args",
			args{key: "1234567890123456789012345678901234567890123456789012345678901234", value: "1234567890123456789012345678901234567890123456789012345678901234"},
			want{err: false},
		},
		{
			"empty args",
			args{key: "", value: ""},
			want{err: false},
		},
		{
			"storage return error",
			args{key: "", value: ""},
			want{err: true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupTest(t)
			if tt.want.err {
				storageMock.
					EXPECT().
					Put(mock.AnythingOfType("string"), mock.AnythingOfType("string")).
					Return(errors.New("")).Times(1)
			} else {
				storageMock.
					EXPECT().
					Put(tt.args.key, tt.args.value).
					Return(nil).
					Times(1)
				tLoggerMock.
					EXPECT().
					WritePut(tt.args.key, tt.args.value).
					Times(1)
			}

			err := srv.Put(tt.args.key, tt.args.value)

			assert.Equal(t, tt.want.err, err != nil)
		})
	}
}

func TestKeyService_Get(t *testing.T) {
	type args struct {
		key string
	}
	type want struct {
		err   bool
		value string
	}
	type test struct {
		name string
		args args
		want want
	}
	tests := []test{
		{
			"get existing value",
			args{key: "1"},
			want{value: "one", err: false},
		},
		{
			"storage return error",
			args{key: "1"},
			want{value: "", err: true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupTest(t)
			if tt.want.err {
				storageMock.
					EXPECT().
					Get(mock.AnythingOfType("string")).
					Return("", errors.New("")).
					Times(1)
			} else {
				storageMock.
					EXPECT().
					Get(tt.args.key).
					Return(tt.want.value, nil).
					Times(1)
			}

			value, err := srv.Get(tt.args.key)

			assert.Equal(t, tt.want.err, err != nil)
			assert.Equal(t, tt.want.value, value)
		})
	}
}

func TestKeyService_Delete(t *testing.T) {
	type args struct {
		key string
	}
	type want struct {
		err bool
	}
	type test struct {
		name string
		args args
		want want
	}
	tests := []test{
		{
			"delete existing value",
			args{key: "1"},
			want{err: false},
		},
		{
			"storage return error",
			args{key: "1"},
			want{err: true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupTest(t)
			if tt.want.err {
				storageMock.
					EXPECT().
					Delete(mock.AnythingOfType("string")).
					Return(errors.New("")).
					Times(1)
			} else {
				storageMock.
					EXPECT().
					Delete(mock.AnythingOfType("string")).
					Return(nil).
					Times(1)
				tLoggerMock.
					EXPECT().
					WriteDelete(mock.AnythingOfType("string")).
					Times(1)
			}

			err := srv.Delete(tt.args.key)

			assert.Equal(t, tt.want.err, err != nil)
		})
	}
}
