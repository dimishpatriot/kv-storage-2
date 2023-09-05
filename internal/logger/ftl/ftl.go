package ftl

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/dimishpatriot/kv-storage/internal/logger"
)

type FileTransactionLogger struct {
	events       chan<- logger.Event
	errors       <-chan error
	lastSequence uint64
	file         *os.File
}

func NewFileTransactionLogger(filename string) (logger.TransactionLogger, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0o755)
	if err != nil {
		return nil, fmt.Errorf("cannot open log file: %w", err)
	}

	return &FileTransactionLogger{file: file}, nil
}

func (l *FileTransactionLogger) Run() {
	events := make(chan logger.Event, 16)
	l.events = events

	errors := make(chan error, 1)
	l.errors = errors

	go func() {
		defer l.file.Close()

		for e := range events {
			l.lastSequence++
			_, err := fmt.Fprintf(
				l.file,
				"%d\t%d\t%s\t%s\n",
				l.lastSequence, e.EventType, e.Key, e.Value,
			)
			if err != nil {
				errors <- err
				return
			}
		}
	}()
}

func (l *FileTransactionLogger) ReadEvents() (<-chan logger.Event, <-chan error) {
	scanner := bufio.NewScanner(l.file)
	outEvent := make(chan logger.Event)
	outError := make(chan error, 1)

	go func() {
		var e logger.Event
		defer close(outEvent)
		defer close(outError)

		for scanner.Scan() {
			line := scanner.Text()

			_, err := fmt.Sscanf(line, "%d\t%d\t%s\t%s", &e.Sequence, &e.EventType, &e.Key, &e.Value)
			if err != nil && !errors.Is(err, io.EOF) {
				outError <- fmt.Errorf("input parse error: %w", err)
				return
			}

			if l.lastSequence >= e.Sequence {
				outError <- fmt.Errorf("transaction numbers out of sequence")
				return
			}

			l.lastSequence = e.Sequence
			outEvent <- e
		}

		if err := scanner.Err(); err != nil {
			outError <- fmt.Errorf("transaction log read failure: %w", err)
			return
		}
	}()

	return outEvent, outError
}

func (l *FileTransactionLogger) WritePut(key, value string) {
	l.events <- logger.Event{
		EventType: logger.EventPut, Key: key, Value: value,
	}
}

func (l *FileTransactionLogger) WriteDelete(key string) {
	l.events <- logger.Event{
		EventType: logger.EventDelete, Key: key,
	}
}

func (l *FileTransactionLogger) Err() <-chan error {
	return l.errors
}
