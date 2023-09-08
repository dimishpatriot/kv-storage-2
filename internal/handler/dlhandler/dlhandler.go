package dlhandler

import (
	"errors"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/dimishpatriot/kv-storage/internal/handler"
	"github.com/dimishpatriot/kv-storage/internal/storage"
	"github.com/dimishpatriot/kv-storage/internal/transactionlogger"
	"github.com/gorilla/mux"
)

type DataLoggerHandler struct {
	logger     *log.Logger
	dataLogger transactionlogger.TransactionLogger
	storage    storage.Storage
}

func New(logger *log.Logger, dataLogger transactionlogger.TransactionLogger, storage storage.Storage) handler.Handler {
	return &DataLoggerHandler{
		dataLogger: dataLogger,
		logger:     logger,
		storage:    storage,
	}
}

func (dh *DataLoggerHandler) Put(w http.ResponseWriter, r *http.Request) {
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

	err = dh.storage.Put(key, value)
	if err != nil {
		http.Error(w,
			err.Error(),
			http.StatusInternalServerError)
		return
	}

	dh.logger.Printf("put: {%s: %s}\n", key, value)
	dh.dataLogger.WritePut(key, value)
	w.WriteHeader(http.StatusCreated)
}

func (dh *DataLoggerHandler) Get(w http.ResponseWriter, r *http.Request) {
	var err error

	key := mux.Vars(r)["key"]
	err = checkKey(key)
	if err != nil {
		http.Error(w,
			err.Error(),
			http.StatusBadRequest)
	}

	value, err := dh.storage.Get(key)
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

	dh.logger.Printf("get: {%s: %s}\n", key, value)
	_, _ = w.Write([]byte(value))
}

func (dh *DataLoggerHandler) Delete(w http.ResponseWriter, r *http.Request) {
	var err error

	key := mux.Vars(r)["key"]
	err = checkKey(key)
	if err != nil {
		http.Error(w,
			err.Error(),
			http.StatusBadRequest)
	}

	err = dh.storage.Delete(key)
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

	dh.logger.Printf("delete: {%s}\n", key)
	dh.dataLogger.WriteDelete(key)
	w.WriteHeader(http.StatusOK)
}

func checkKey(key string) error {
	if key == "" {
		return handler.ErrorEmptyKey
	}
	if len(key) > 64 {
		return handler.ErrorLongKey
	}
	if strings.ContainsAny(key, " /\t\n") {
		return handler.ErrorKeyContainsForbiddenSymbol
	}

	return nil
}

func checkValue(value string) error {
	if value == "" {
		return handler.ErrorEmptyValue
	}
	if len(value) > 128 {
		return handler.ErrorLongValue
	}

	return nil
}
