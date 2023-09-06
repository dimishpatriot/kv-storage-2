package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/dimishpatriot/kv-storage/internal/datalogger"
	"github.com/dimishpatriot/kv-storage/internal/handler"
	"github.com/dimishpatriot/kv-storage/internal/handler/dlhandler"
	"github.com/dimishpatriot/kv-storage/internal/storage"
	"github.com/dimishpatriot/kv-storage/internal/storage/localstorage"
	"github.com/gorilla/mux"
)

func main() {
	app, err := New()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	log.Fatal(app.Run())
}

type App struct {
	logger     *log.Logger
	router     *mux.Router
	dataLogger datalogger.TransactionLogger
	handler    handler.Handler
	storage    storage.Storage
}

func New() (*App, error) {
	logger := log.New(os.Stdout, "INFO:", log.Lshortfile|log.Ltime|log.Lmicroseconds|log.Ldate)
	logger.Println("logger created")

	storage := localstorage.New()
	logger.Println("storage created")

	dataLogger, err := datalogger.New(logger, "transaction.log", storage)
	if err != nil {
		return nil, fmt.Errorf("failed to create data logger: %w", err)
	}
	logger.Println("dataLogger created")

	handler := dlhandler.New(logger, dataLogger, storage)
	logger.Println("controller created")

	router := mux.NewRouter()
	logger.Println("router created")

	return &App{
		dataLogger: dataLogger,
		handler:    handler,
		logger:     logger,
		router:     router,
		storage:    storage,
	}, nil
}

func (a *App) addRoutes() {
	a.router.HandleFunc("/v1/{key}", a.handler.Put).Methods("PUT")
	a.router.HandleFunc("/v1/{key}", a.handler.Get).Methods("GET")
	a.router.HandleFunc("/v1/{key}", a.handler.Delete).Methods("DELETE")
}

func (a *App) Run() error {
	a.addRoutes()
	a.logger.Println("routes added")

	a.dataLogger.RestoreDataFromFile()
	a.logger.Println("data restored")

	a.dataLogger.Run()
	a.logger.Println("dataLogger ran")

	return http.ListenAndServe(":8080", a.router)
}
