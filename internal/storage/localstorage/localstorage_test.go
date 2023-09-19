package localstorage_test

import (
	"errors"
	"testing"

	"github.com/dimishpatriot/kv-storage/internal/storage"
	"github.com/dimishpatriot/kv-storage/internal/storage/localstorage"
)

var store *localstorage.LocalStorage

func setupTest(tb testing.TB) func(tb testing.TB) {
	store = localstorage.New().(*localstorage.LocalStorage)
	_ = store.Put("one", "ONE")
	_ = store.Put("0123456789", "numbers")

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
		{
			"correct key",
			args{key: "correct key", value: "correct value"},
			nil,
		},
		{
			"existing key",
			args{key: "one", value: "NEW ONE"},
			nil,
		},
		{
			"one symbol key",
			args{key: "K", value: "correct value"},
			nil,
		},
		{
			"non alphabet key",
			args{key: "~!@#$%^&*()_+", value: "correct value"},
			nil,
		},
		{
			"number string key",
			args{key: "0123456789", value: "correct value"},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupTest(t)

			err := store.Put(tt.args.key, tt.args.value)
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
		{
			"common key",
			args{key: "one"},
			want{value: "ONE", err: nil},
		},
		{
			"number string key",
			args{key: "0123456789"},
			want{value: "numbers", err: nil},
		},
		{
			"absent key",
			args{key: "absent"},
			want{value: "", err: storage.ErrorNoSuchKey},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupTest(t)

			got, err := store.Get(tt.args.key)
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
		{
			"exists key",
			args{key: "one"},
			nil,
		},
		{
			"absent key",
			args{key: "ONE"},
			storage.ErrorNoSuchKey,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupTest(t)

			err := store.Delete(tt.args.key)
			if (err != nil) && !errors.Is(err, tt.wantErr) {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
