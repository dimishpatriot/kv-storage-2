package controller

import (
	"errors"
	"io"
	"log"
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
	logger     *log.Logger
	dataLogger datalogger.TransactionLogger
	Storage    storage.Storage
}

func New(logger *log.Logger, dataLogger datalogger.TransactionLogger) *Controller {
	s := localstorage.New()
	return &Controller{dataLogger: dataLogger, Storage: s, logger: logger}
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

	bValue, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w,
			err.Error(),
			http.StatusInternalServerError)
		return
	}

	value := string(bValue)

	err = checkValue(value)
	if err != nil {
		http.Error(w,
			err.Error(),
			http.StatusBadRequest)
	}

	err = c.Storage.Put(key, value)
	if err != nil {
		http.Error(w,
			err.Error(),
			http.StatusInternalServerError)
		return
	}

	c.logger.Printf("put: {%s: %s}\n", key, value)
	c.dataLogger.WritePut(key, value)
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

	c.logger.Printf("get: {%s: %s}\n", key, value)
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

	c.logger.Printf("delete: {%s}\n", key)
	c.dataLogger.WriteDelete(key)
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
