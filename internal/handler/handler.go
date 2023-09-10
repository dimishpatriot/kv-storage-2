package handler

import (
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/dimishpatriot/kv-storage/internal/services/keyservice"
	"github.com/dimishpatriot/kv-storage/internal/storage"
	"github.com/gorilla/mux"
)

//go:generate mockery --name Handler
type Handler interface {
	Put(http.ResponseWriter, *http.Request)
	Get(http.ResponseWriter, *http.Request)
	Delete(http.ResponseWriter, *http.Request)
}

type dataHandler struct {
	keyService keyservice.KeyService
}

var (
	ErrorEmptyKey                   = errors.New("empty key")
	ErrorLongKey                    = errors.New("key length > 64 byte")
	ErrorKeyContainsForbiddenSymbol = errors.New("forbidden symbol in key")
	ErrorEmptyValue                 = errors.New("empty value")
	ErrorLongValue                  = errors.New("value length > 128 byte")
)

func New(keyService keyservice.KeyService) Handler {
	return &dataHandler{keyService}
}

func (dh *dataHandler) Put(w http.ResponseWriter, r *http.Request) {
	key, err := dh.getKeyFromRequest(r)
	if err != nil {
		http.Error(w,
			err.Error(),
			http.StatusBadRequest)
		return
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
		return
	}

	err = dh.keyService.Put(key, value)
	if err != nil {
		http.Error(w,
			err.Error(),
			http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (dh *dataHandler) Get(w http.ResponseWriter, r *http.Request) {
	key, err := dh.getKeyFromRequest(r)
	if err != nil {
		http.Error(w,
			err.Error(),
			http.StatusBadRequest)
		return
	}

	value, err := dh.keyService.Get(key)
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

func (dh *dataHandler) Delete(w http.ResponseWriter, r *http.Request) {
	key, err := dh.getKeyFromRequest(r)
	if err != nil {
		http.Error(w,
			err.Error(),
			http.StatusBadRequest)
		return
	}

	err = dh.keyService.Delete(key)
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

	w.WriteHeader(http.StatusOK)
}

func (dh *dataHandler) getKeyFromRequest(r *http.Request) (string, error) {
	key := mux.Vars(r)["key"]
	return key, checkKey(key)
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
