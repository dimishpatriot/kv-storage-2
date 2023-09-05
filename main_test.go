package main

import (
	"fmt"
	"testing"
)

func setupTest(tb testing.TB) func(tb testing.TB) {
	// run before test
	store = map[string]string{
		"one":        "ONE",
		"0123456789": "numbers",
	}

	return func(tb testing.TB) {
		// run after test
		// ...
	}
}

func TestPut(t *testing.T) {
	type args struct {
		k string
		v string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "correct key", args: args{k: "correct key", v: "correct value"}, wantErr: false},
		{name: "existing key", args: args{k: "one", v: "NEW ONE"}, wantErr: false},
		{name: "one symbol key", args: args{k: "K", v: "correct value"}, wantErr: false},
		{name: "non alphabet key", args: args{k: "/~!@#$%^&*()_+", v: "correct value"}, wantErr: false},
		{name: "number string key", args: args{k: "0123456789", v: "correct value"}, wantErr: false},
		{name: "empty key", args: args{k: "", v: ""}, wantErr: true},
		{name: "empty value", args: args{k: "", v: ""}, wantErr: true},
		{name: "empty key & value", args: args{k: "", v: ""}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			after := setupTest(t)
			defer after(t)

			if err := Put(tt.args.k, tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("Put() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	fmt.Println(store)
}

func TestGet(t *testing.T) {
	type args struct {
		k string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{name: "common key", args: args{k: "one"}, want: "ONE", wantErr: false},
		{name: "number string key", args: args{k: "0123456789"}, want: "numbers", wantErr: false},
		{name: "empty key", args: args{k: ""}, want: "", wantErr: true},
		{name: "absent key", args: args{k: "absent"}, want: "", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			after := setupTest(t)
			defer after(t)

			got, err := Get(tt.args.k)
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
		k string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "exists key", args: args{k: "one"}, wantErr: false},
		{name: "absent key", args: args{k: "ONE"}, wantErr: true},
		{name: "empty key", args: args{k: ""}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			after := setupTest(t)
			defer after(t)

			if err := Delete(tt.args.k); (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
