package main

import (
	"errors"
)

var store = make(map[string]string)

var (
	ErrorEmptyData = errors.New("empty data")
	ErrorNoSuchKey = errors.New("no such key")
)

func Put(k string, v string) error {
	if k == "" || v == "" {
		return ErrorEmptyData
	}
	store[k] = v

	return nil
}

func Get(k string) (string, error) {
	v, ok := store[k]
	if !ok {
		return "", ErrorNoSuchKey
	}

	return v, nil
}

func Delete(k string) error {
	if _, ok := store[k]; k == "" || !ok {
		return ErrorNoSuchKey
	}
	delete(store, k)

	return nil
}

func main() {
}
