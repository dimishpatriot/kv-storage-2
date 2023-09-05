package main

import (
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/dimishpatriot/kv-storage/internal/service"
	"github.com/gorilla/mux"
)

func PutHandler(w http.ResponseWriter, r *http.Request) {
	var err error

	key := mux.Vars(r)["key"]
	err = checkKey(key)
	if err != nil {
		http.Error(w,
			err.Error(),
			http.StatusBadRequest)
	}

	value, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w,
			err.Error(),
			http.StatusInternalServerError)
		return
	}
	err = checkValue(string(value))
	if err != nil {
		http.Error(w,
			err.Error(),
			http.StatusBadRequest)
	}

	err = service.Put(key, string(value))
	if err != nil {
		http.Error(w,
			err.Error(),
			http.StatusInternalServerError)
		return
	}

	l.WritePut(key, string(value))
	w.WriteHeader(http.StatusCreated)
}

func GetHandler(w http.ResponseWriter, r *http.Request) {
	var err error

	key := mux.Vars(r)["key"]
	err = checkKey(key)
	if err != nil {
		http.Error(w,
			err.Error(),
			http.StatusBadRequest)
	}

	value, err := service.Get(key)
	if errors.Is(err, service.ErrorNoSuchKey) {
		http.Error(w,
			err.Error(),
			http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w,
			err.Error(),
			http.StatusInternalServerError)
		return
	}

	_, _ = w.Write([]byte(value))
}

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	var err error

	key := mux.Vars(r)["key"]
	err = checkKey(key)
	if err != nil {
		http.Error(w,
			err.Error(),
			http.StatusBadRequest)
	}

	err = service.Delete(key)
	if errors.Is(err, service.ErrorNoSuchKey) {
		http.Error(w,
			err.Error(),
			http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w,
			err.Error(),
			http.StatusInternalServerError)
		return
	}

	l.WriteDelete(key)
	w.WriteHeader(http.StatusOK)
}

var (
	ErrorEmptyKey                   = errors.New("empty key")
	ErrorLongKey                    = errors.New("key length > 64 byte")
	ErrorKeyContainsForbiddenSymbol = errors.New("forbidden symbol in key")
	ErrorEmptyValue                 = errors.New("empty value")
	ErrorLongValue                  = errors.New("value length > 128 byte")
)

func checkKey(key string) error {
	if key == "" {
		return ErrorEmptyKey
	}
	if len(key) > 64 {
		return ErrorLongKey
	}
	if strings.ContainsAny(key, " /\t\n") {
		return ErrorKeyContainsForbiddenSymbol
	}

	return nil
}

func checkValue(value string) error {
	if value == "" {
		return ErrorEmptyValue
	}
	if len(value) > 128 {
		return ErrorLongValue
	}

	return nil
}
