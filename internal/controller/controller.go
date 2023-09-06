package controller

import (
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/dimishpatriot/kv-storage/internal/logger/datalogger"
	"github.com/dimishpatriot/kv-storage/internal/storage"
	"github.com/dimishpatriot/kv-storage/internal/storage/localstorage"
	"github.com/gorilla/mux"
)

var (
	ErrorEmptyKey                   = errors.New("empty key")
	ErrorLongKey                    = errors.New("key length > 64 byte")
	ErrorKeyContainsForbiddenSymbol = errors.New("forbidden symbol in key")
	ErrorEmptyValue                 = errors.New("empty value")
	ErrorLongValue                  = errors.New("value length > 128 byte")
)

type Controller struct {
	logger  datalogger.TransactionLogger
	Storage storage.Storage
}

func New(logger datalogger.TransactionLogger) *Controller {
	s := localstorage.New()
	return &Controller{logger: logger, Storage: s}
}

func (c *Controller) PutHandler(w http.ResponseWriter, r *http.Request) {
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

	err = c.Storage.Put(key, string(value))
	if err != nil {
		http.Error(w,
			err.Error(),
			http.StatusInternalServerError)
		return
	}

	c.logger.WritePut(key, string(value))
	w.WriteHeader(http.StatusCreated)
}

func (c *Controller) GetHandler(w http.ResponseWriter, r *http.Request) {
	var err error

	key := mux.Vars(r)["key"]
	err = checkKey(key)
	if err != nil {
		http.Error(w,
			err.Error(),
			http.StatusBadRequest)
	}

	value, err := c.Storage.Get(key)
	if errors.Is(err, storage.ErrorNoSuchKey) {
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

func (c *Controller) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	var err error

	key := mux.Vars(r)["key"]
	err = checkKey(key)
	if err != nil {
		http.Error(w,
			err.Error(),
			http.StatusBadRequest)
	}

	err = c.Storage.Delete(key)
	if errors.Is(err, storage.ErrorNoSuchKey) {
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

	c.logger.WriteDelete(key)
	w.WriteHeader(http.StatusOK)
}

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
