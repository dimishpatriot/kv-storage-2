package localstorage

import (
	"sync"

	"github.com/dimishpatriot/kv-storage/internal/storage"
)

type data = map[string]string

type LocalStorage struct {
	sync.RWMutex
	data data
}

func New() storage.Storage {
	d := make(map[string]string)

	return &LocalStorage{data: d}
}

func (ls *LocalStorage) set(data data) {
	ls.Lock()
	ls.data = data
	ls.Unlock()
}

func (ls *LocalStorage) Put(k string, v string) error {
	ls.Lock()
	ls.data[k] = v
	ls.Unlock()

	return nil
}

func (ls *LocalStorage) Get(k string) (string, error) {
	ls.RLock()
	v, ok := ls.data[k]
	ls.RUnlock()
	if !ok {
		return "", storage.ErrorNoSuchKey
	}

	return v, nil
}

func (ls *LocalStorage) Delete(k string) error {
	ls.Lock()
	defer ls.Unlock()
	if _, ok := ls.data[k]; !ok {
		return storage.ErrorNoSuchKey
	}
	delete(ls.data, k)

	return nil
}
