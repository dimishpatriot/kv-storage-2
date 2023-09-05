package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dimishpatriot/kv-storage/internal/logger"
	"github.com/dimishpatriot/kv-storage/internal/logger/ftl"
	"github.com/dimishpatriot/kv-storage/internal/service"
	"github.com/gorilla/mux"
)

var l logger.TransactionLogger

func main() {
	err := initTransactionLog()
	if err != nil {
		log.Fatalf("cannot initialize logger: %s", err)
	}

	r := mux.NewRouter()

	r.HandleFunc("/v1/{key}", PutHandler).Methods("PUT")
	r.HandleFunc("/v1/{key}", GetHandler).Methods("GET")
	r.HandleFunc("/v1/{key}", DeleteHandler).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", r))
}

func initTransactionLog() error {
	var err error

	l, err = ftl.NewFileTransactionLogger("transaction.log")
	if err != nil {
		return fmt.Errorf("failed to create event logger: %w", err)
	}

	events, errors := l.ReadEvents()
	e, ok := logger.Event{}, true

	for ok && err == nil {
		select {
		case err, ok = <-errors:
		case e, ok = <-events:
			switch e.EventType {
			case logger.EventDelete:
				err = service.Delete(e.Key)
			case logger.EventPut:
				err = service.Put(e.Key, e.Value)
			}
		}
	}

	l.Run()

	return err
}
