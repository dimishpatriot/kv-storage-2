package service

import (
	"errors"
	"sync"
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
