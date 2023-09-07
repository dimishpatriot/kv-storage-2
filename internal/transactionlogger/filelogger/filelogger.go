package filelogger

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/dimishpatriot/kv-storage/internal/storage"
	"github.com/dimishpatriot/kv-storage/internal/transactionlogger"
)

const (
	readPattern  = "%d\t%d\t%s\t%s"
	writePattern = "%d\t%d\t%s\t%s\n"
)

type FileTransactionLogger struct {
	events       chan<- transactionlogger.Event
	errors       <-chan error
	lastSequence uint64
	file         *os.File
	logger       *log.Logger
	storage      storage.Storage
}

func New(logger *log.Logger, filename string, storage storage.Storage) (transactionlogger.TransactionLogger, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0o755)
	if err != nil {
		return nil, fmt.Errorf("cant open log file: %w", err)
	}
	ftl := FileTransactionLogger{logger: logger, file: file, storage: storage}
	ftl.restoreData()

	return &ftl, nil
}

func (l *FileTransactionLogger) restoreData() {
	l.logger.Println("restore data...")

	var err error
	events, errors := l.ReadEvents()
	e, ok := transactionlogger.Event{}, true

	for ok && err == nil {
		select {
		case err, ok = <-errors:
		case e, ok = <-events:
			switch e.EventType {
			case transactionlogger.EventDelete:
				err = l.storage.Delete(e.Key)
			case transactionlogger.EventPut:
				err = l.storage.Put(e.Key, e.Value)
			}
		}
	}
}

func (l *FileTransactionLogger) Run() {
	l.logger.Println("dataLogger run...")

	events := make(chan transactionlogger.Event, 16)
	l.events = events
	errors := make(chan error, 1)
	l.errors = errors

	go func() {
		defer l.file.Close()

		for e := range events {
			l.lastSequence++
			_, err := fmt.Fprintf(l.file, writePattern, l.lastSequence, e.EventType, e.Key, e.Value)
			if err != nil {
				errors <- err
				return
			}
			if e.EventType == transactionlogger.EventDelete {
				if err = l.clearNotActualData(e.Key); err != nil {
					errors <- err
					return
				}
			}
		}
	}()
}

func (l *FileTransactionLogger) clearNotActualData(key string) error {
	l.logger.Println("clear not actual data...")

	tempFileName := "temp.log"
	tempFile, err := os.OpenFile(tempFileName, os.O_WRONLY|os.O_CREATE, 0o755)
	if err != nil {
		return fmt.Errorf("cant create temp log file: %w", err)
	}
	defer tempFile.Close()

	_, _ = l.file.Seek(0, 0) // seek to start!
	if err = l.copyData(key, tempFile); err != nil {
		return fmt.Errorf("cant copy data: %w", err)
	}

	l.file.Close()
	if err = l.swapFiles(l.file.Name(), tempFileName); err != nil {
		return fmt.Errorf("cant swap log files: %w", err)
	}

	l.file, err = os.OpenFile(l.file.Name(), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o755)
	if err != nil {
		return fmt.Errorf("cant open new log file: %w", err)
	}
	_, _ = l.file.Seek(0, 2) // seek to end!
	return nil
}

func (l *FileTransactionLogger) swapFiles(logFileName string, tempFileName string) error {
	l.logger.Println("swap log files...")

	if err := os.Remove(logFileName); err != nil {
		return fmt.Errorf("cant remove old log file: %w", err)
	}
	if err := os.Rename(tempFileName, logFileName); err != nil {
		return fmt.Errorf("cant rename temp log file: %w", err)
	}
	return nil
}

func (l *FileTransactionLogger) copyData(key string, tempFile *os.File) error {
	l.logger.Println("coping log data...")

	var err error
	var e transactionlogger.Event

	scanner := bufio.NewScanner(l.file)
	for scanner.Scan() {
		line := scanner.Text()

		_, err = fmt.Sscanf(line, readPattern, &e.Sequence, &e.EventType, &e.Key, &e.Value)
		if err != nil && !errors.Is(err, io.EOF) {
			return fmt.Errorf("input parse error: %w", err)
		}
		if e.Key != key {
			_, err = fmt.Fprintf(tempFile, writePattern, e.Sequence, e.EventType, e.Key, e.Value)
			if err != nil {
				return fmt.Errorf("cant save to temp file: %w", err)
			}
		}
	}
	return nil
}

func (l *FileTransactionLogger) ReadEvents() (<-chan transactionlogger.Event, <-chan error) {
	l.logger.Println("read events...")

	scanner := bufio.NewScanner(l.file)
	outEvent := make(chan transactionlogger.Event)
	outError := make(chan error, 1)

	go func() {
		var e transactionlogger.Event
		defer close(outEvent)
		defer close(outError)

		for scanner.Scan() {
			line := scanner.Text()

			_, err := fmt.Sscanf(line, readPattern, &e.Sequence, &e.EventType, &e.Key, &e.Value)
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
	l.logger.Printf("write put: {%s: %s}", key, value)

	l.events <- transactionlogger.Event{
		EventType: transactionlogger.EventPut, Key: key, Value: value,
	}
}

func (l *FileTransactionLogger) WriteDelete(key string) {
	l.logger.Printf("write delete {%s}", key)

	l.events <- transactionlogger.Event{
		EventType: transactionlogger.EventDelete, Key: key,
	}
}

func (l *FileTransactionLogger) Err() <-chan error {
	l.logger.Println("getting errors channel...")

	return l.errors
}
