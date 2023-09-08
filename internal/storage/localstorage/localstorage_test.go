package localstorage_test

import (
	"testing"

	"github.com/dimishpatriot/kv-storage/internal/storage/localstorage"
)

var store *localstorage.LocalStorage

func setupTest(tb testing.TB) func(tb testing.TB) {
	// run before each test
	store = localstorage.New().(*localstorage.LocalStorage)
	_ = store.Put("one", "ONE")
	_ = store.Put("0123456789", "numbers")

	return func(tb testing.TB) {
		// run after each test
		// ...
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
		wantErr bool
	}{
		{name: "correct key", args: args{key: "correct key", value: "correct value"}, wantErr: false},
		{name: "existing key", args: args{key: "one", value: "NEW ONE"}, wantErr: false},
		{name: "one symbol key", args: args{key: "K", value: "correct value"}, wantErr: false},
		{name: "non alphabet key", args: args{key: "~!@#$%^&*()_+", value: "correct value"}, wantErr: false},
		{name: "number string key", args: args{key: "0123456789", value: "correct value"}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			after := setupTest(t)
			defer after(t)

			if err := store.Put(tt.args.key, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("Put() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGet(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{name: "common key", args: args{key: "one"}, want: "ONE", wantErr: false},
		{name: "number string key", args: args{key: "0123456789"}, want: "numbers", wantErr: false},
		{name: "absent key", args: args{key: "absent"}, want: "", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			after := setupTest(t)
			defer after(t)

			got, err := store.Get(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Get() = %v, want %v", got, tt.want)
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
		wantErr bool
	}{
		{name: "exists key", args: args{key: "one"}, wantErr: false},
		{name: "absent key", args: args{key: "ONE"}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			after := setupTest(t)
			defer after(t)

			if err := store.Delete(tt.args.key); (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
