package handler

import (
	"errors"
	"testing"
)

func TestDataHandler_checkKey(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name          string
		args          args
		wantErr       bool
		wantErrorType error
	}{
		{
			"correct simple key",
			args{key: "123abc"},
			false,
			nil,
		},
		{
			"correct long key",
			args{key: "1234567890123456789012345678901234567890123456789012345678901234"},
			false,
			nil,
		},
		{
			"key with symbols",
			args{key: "!@#$%^&*()_+><"},
			false,
			nil,
		},
		{
			"empty key",
			args{key: ""},
			true,
			ErrorEmptyKey,
		},
		{
			"very long key",
			args{key: "12345678901234567890123456789012345678901234567890123456789012345"},
			true,
			ErrorLongKey,
		},
		{
			"key with space",
			args{key: "abc def"},
			true,
			ErrorKeyContainsForbiddenSymbol,
		},
		{
			"key with tab",
			args{key: "abc\tdef"},
			true,
			ErrorKeyContainsForbiddenSymbol,
		},
		{
			"key with new line",
			args{key: "abcdef\n"},
			true,
			ErrorKeyContainsForbiddenSymbol,
		},
		{
			"key with slash",
			args{key: "a/b"},
			true,
			ErrorKeyContainsForbiddenSymbol,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := checkKey(tt.args.key)

			if (err != nil) != tt.wantErr {
				t.Errorf("checkKey() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && !errors.Is(err, tt.wantErrorType) {
				t.Errorf("checkKey() errorType = %v, want %v", err, tt.wantErrorType)
			}
		})
	}
}

func TestDataHandler_checkValue(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name          string
		args          args
		wantErr       bool
		wantErrorType error
	}{
		{
			"correct simple value",
			args{value: "123abc"},
			false,
			nil,
		},
		{
			"correct long value",
			args{
				"12345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678",
			},
			false,
			nil,
		},
		{
			"value with symbols",
			args{value: "!@#$%^&*()_+><"},
			false,
			nil,
		},
		{
			"multi line value",
			args{value: "abc\ndef\n\n0"},
			false,
			nil,
		},
		{
			"value with tabs",
			args{value: "tab\ttab\t\ttab"},
			false,
			nil,
		},
		{
			"empty value",
			args{value: ""},
			true,
			ErrorEmptyValue,
		},
		{
			"very long value",
			args{
				"123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789",
			},
			true,
			ErrorLongValue,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := checkValue(tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("checkValue() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
