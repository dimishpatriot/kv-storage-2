package localstorage_test

import (
	"errors"
	"testing"

	"github.com/dimishpatriot/kv-storage/internal/storage"
)

var storageMock *storage.MockStorage

func setupTest(tb testing.TB) func(tb testing.TB) {
	storageMock = storage.NewMockStorage(tb)

	return func(tb testing.TB) {
	}
}

func TestPut(t *testing.T) {
	type args struct {
		key   string
		value string
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{name: "correct key", args: args{key: "correct key", value: "correct value"}, wantErr: nil},
		{name: "existing key", args: args{key: "one", value: "NEW ONE"}, wantErr: nil},
		{name: "one symbol key", args: args{key: "K", value: "correct value"}, wantErr: nil},
		{name: "non alphabet key", args: args{key: "~!@#$%^&*()_+", value: "correct value"}, wantErr: nil},
		{name: "number string key", args: args{key: "0123456789", value: "correct value"}, wantErr: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			after := setupTest(t)
			defer after(t)
			storageMock.EXPECT().Put(tt.args.key, tt.args.value).Return(nil)

			err := storageMock.Put(tt.args.key, tt.args.value)
			if (err != nil) && !errors.Is(err, tt.wantErr) {
				t.Errorf("Put() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGet(t *testing.T) {
	type args struct {
		key string
	}
	type want struct {
		value string
		err   error
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{name: "common key", args: args{key: "one"}, want: want{value: "ONE", err: nil}},
		{name: "number string key", args: args{key: "0123456789"}, want: want{value: "numbers", err: nil}},
		{name: "absent key", args: args{key: "absent"}, want: want{value: "", err: storage.ErrorNoSuchKey}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			after := setupTest(t)
			defer after(t)
			storageMock.EXPECT().Get(tt.args.key).Return(tt.want.value, tt.want.err)

			got, err := storageMock.Get(tt.args.key)
			if (err != nil) && !errors.Is(err, tt.want.err) {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.want.err)
				return
			}
			if got != tt.want.value {
				t.Errorf("Get() = %s, want %s", got, tt.want.value)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{name: "exists key", args: args{key: "one"}, wantErr: nil},
		{name: "absent key", args: args{key: "ONE"}, wantErr: storage.ErrorNoSuchKey},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			after := setupTest(t)
			defer after(t)
			storageMock.EXPECT().Delete(tt.args.key).Return(nil)

			err := storageMock.Delete(tt.args.key)
			if (err != nil) && !errors.Is(err, tt.wantErr) {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
