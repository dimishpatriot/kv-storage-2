package localstorage

import (
	"sync"

	"github.com/dimishpatriot/kv-storage/internal/storage"
)

type LocalStorage struct {
	sync.RWMutex
	m map[string]string
}

func New() *LocalStorage {
	m := make(map[string]string)

	return &LocalStorage{m: m}
}

func (ls *LocalStorage) Put(k string, v string) error {
	ls.Lock()
	ls.m[k] = v
	ls.Unlock()

	return nil
}

func (ls *LocalStorage) Get(k string) (string, error) {
	ls.RLock()
	v, ok := ls.m[k]
	ls.RUnlock()
	if !ok {
		return "", storage.ErrorNoSuchKey
	}

	return v, nil
}

func (ls *LocalStorage) Delete(k string) error {
	ls.Lock()
	defer ls.Unlock()
	if _, ok := ls.m[k]; !ok {
		return storage.ErrorNoSuchKey
	}
	delete(ls.m, k)

	return nil
}
