package dlhandler

import (
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/dimishpatriot/kv-storage/internal/handler"
	"github.com/dimishpatriot/kv-storage/internal/services/keyservice"
	"github.com/dimishpatriot/kv-storage/internal/storage"
	"github.com/gorilla/mux"
)

type dataHandler struct {
	keyService keyservice.KeyService
}

func New(keyService keyservice.KeyService) handler.Handler {
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
