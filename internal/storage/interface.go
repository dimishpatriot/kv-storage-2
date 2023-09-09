package storage

import "errors"

//go:generate mockery --name Storage
type Storage interface {
	Put(string, string) error
	Get(string) (string, error)
	Delete(string) error
}

var ErrorNoSuchKey = errors.New("no such key")
