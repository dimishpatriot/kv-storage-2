package handler

import (
	"errors"
	"net/http"
)

//go:generate mockery --name Handler
type Handler interface {
	Put(http.ResponseWriter, *http.Request)
	Get(http.ResponseWriter, *http.Request)
	Delete(http.ResponseWriter, *http.Request)
}

var (
	ErrorEmptyKey                   = errors.New("empty key")
	ErrorLongKey                    = errors.New("key length > 64 byte")
	ErrorKeyContainsForbiddenSymbol = errors.New("forbidden symbol in key")
	ErrorEmptyValue                 = errors.New("empty value")
	ErrorLongValue                  = errors.New("value length > 128 byte")
)
