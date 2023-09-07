package transactionlogger

type TransactionLogger interface {
	Err() <-chan error
	ReadEvents() (<-chan Event, <-chan error)
	Run()
	WriteDelete(key string)
	WritePut(key, value string)
}

type Event struct {
	Sequence  uint64
	EventType EventType
	Key       string
	Value     string
}

type EventType byte

const (
	_                     = iota
	EventDelete EventType = iota
	EventPut
)
