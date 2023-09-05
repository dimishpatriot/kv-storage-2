package main

import (
	"errors"
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

var store = struct {
	sync.RWMutex
	m map[string]string
}{
	m: make(map[string]string),
}

var (
	ErrorEmptyData = errors.New("empty data")
	ErrorNoSuchKey = errors.New("no such key")
)

func Put(k string, v string) error {
	if k == "" || v == "" {
		return ErrorEmptyData
	}
	store.Lock()
	store.m[k] = v
	store.Unlock()

	return nil
}

func Get(k string) (string, error) {
	store.RLock()
	v, ok := store.m[k]
	store.RUnlock()
	if !ok {
		return "", ErrorNoSuchKey
	}

	return v, nil
}

func Delete(k string) error {
	store.Lock()
	defer store.Unlock()
	if _, ok := store.m[k]; k == "" || !ok {
		return ErrorNoSuchKey
	}
	delete(store.m, k)

	return nil
}

func kvPutHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	k := vars["key"]

	v, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w,
			err.Error(),
			http.StatusInternalServerError)
		return
	}

	err = Put(k, string(v))
	if err != nil {
		http.Error(w,
			err.Error(),
			http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func kvGetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	k := vars["key"]

	v, err := Get(k)
	if errors.Is(err, ErrorNoSuchKey) {
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

	_, _ = w.Write([]byte(v))
}

func kvDeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	k := vars["key"]

	err := Delete(k)
	if errors.Is(err, ErrorNoSuchKey) {
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

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/v1/{key}", kvPutHandler).Methods("PUT")
	r.HandleFunc("/v1/{key}", kvGetHandler).Methods("GET")
	r.HandleFunc("/v1/{key}", kvDeleteHandler).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", r))
}
