package mocks

import "github.com/dimishpatriot/kv-storage/internal/transactionlogger"

type MockDataLogger struct{}

func (md *MockDataLogger) Err() <-chan error {
	e := make(<-chan error)
	return e
}

func (md *MockDataLogger) ReadEvents() (<-chan transactionlogger.Event, <-chan error) {
	events := make(<-chan transactionlogger.Event)
	e := make(<-chan error)
	return events, e
}

func (md *MockDataLogger) Run() {}

func (md *MockDataLogger) WriteDelete(key string) {}

func (md *MockDataLogger) WritePut(key, value string) {}
