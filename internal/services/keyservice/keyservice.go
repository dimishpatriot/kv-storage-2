package keyservice

import (
	"log"

	"github.com/dimishpatriot/kv-storage/internal/services/transactionlogger"
	"github.com/dimishpatriot/kv-storage/internal/storage"
)

//go:generate mockery --name KeyService
type KeyService interface {
	Put(string, string) error
	Get(string) (string, error)
	Delete(string) error
}

type keyService struct {
	logger  *log.Logger
	storage storage.Storage
	tLogger transactionlogger.TransactionLogger
}

func New(
	logger *log.Logger,
	storage storage.Storage,
	tLogger transactionlogger.TransactionLogger,
) KeyService {
	return &keyService{
		logger,
		storage,
		tLogger,
	}
}

// Put implements Service.
func (s *keyService) Put(k, v string) error {
	err := s.storage.Put(k, v)
	if err == nil {
		s.logger.Printf("put: {%s: %s}\n", k, v)
		s.tLogger.WritePut(k, v)
	}

	return err
}

// Delete implements Service.
func (s *keyService) Delete(k string) error {
	err := s.storage.Delete(k)
	if err == nil {
		s.logger.Printf("delete: {%s}\n", k)
		s.tLogger.WriteDelete(k)
	}

	return err
}

// Get implements Service.
func (s *keyService) Get(k string) (string, error) {
	v, err := s.storage.Get(k)
	if err == nil {
		s.logger.Printf("get: {%s: %s}\n", k, v)
	}

	return v, err
}
