package dlhandler

import (
	"errors"
	"testing"

	"github.com/dimishpatriot/kv-storage/internal/handler"
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
			name:    "correct simple key",
			args:    args{key: "123abc"},
			wantErr: false,
		},
		{
			name:    "correct long key",
			args:    args{key: "1234567890123456789012345678901234567890123456789012345678901234"},
			wantErr: false,
		},
		{name: "key with symbols", args: args{key: "!@#$%^&*()_+><"}, wantErr: false},
		{
			name:          "empty key",
			args:          args{key: ""},
			wantErr:       true,
			wantErrorType: handler.ErrorEmptyKey,
		},
		{
			name: "very long key",
			args: args{
				key: "12345678901234567890123456789012345678901234567890123456789012345",
			},
			wantErr:       true,
			wantErrorType: handler.ErrorLongKey,
		},
		{
			name:          "key with space",
			args:          args{key: "abc def"},
			wantErr:       true,
			wantErrorType: handler.ErrorKeyContainsForbiddenSymbol,
		},
		{
			name:          "key with tab",
			args:          args{key: "abc\tdef"},
			wantErr:       true,
			wantErrorType: handler.ErrorKeyContainsForbiddenSymbol,
		},
		{
			name:          "key with new line",
			args:          args{key: "abcdef\n"},
			wantErr:       true,
			wantErrorType: handler.ErrorKeyContainsForbiddenSymbol,
		},
		{
			name:          "key with slash",
			args:          args{key: "a/b"},
			wantErr:       true,
			wantErrorType: handler.ErrorKeyContainsForbiddenSymbol,
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
			name:    "correct simple value",
			args:    args{value: "123abc"},
			wantErr: false,
		},
		{
			name: "correct long value",
			args: args{
				value: "12345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678",
			},
			wantErr: false,
		},
		{
			name:    "value with symbols",
			args:    args{value: "!@#$%^&*()_+><"},
			wantErr: false,
		},
		{
			name:    "multi line value",
			args:    args{value: "abc\ndef\n\n0"},
			wantErr: false,
		},
		{
			name:    "value with tabs",
			args:    args{value: "tab\ttab\t\ttab"},
			wantErr: false,
		},
		{
			name:          "empty value",
			args:          args{value: ""},
			wantErr:       true,
			wantErrorType: handler.ErrorEmptyValue,
		},
		{
			name: "very long value",
			args: args{
				value: "123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789",
			},
			wantErr:       true,
			wantErrorType: handler.ErrorLongValue,
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
