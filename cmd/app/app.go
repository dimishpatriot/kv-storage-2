package app

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/dimishpatriot/kv-storage/internal/handler"
	"github.com/dimishpatriot/kv-storage/internal/services/keyservice"
	"github.com/dimishpatriot/kv-storage/internal/services/transactionlogger"
	"github.com/dimishpatriot/kv-storage/internal/services/transactionlogger/filelogger"
	"github.com/dimishpatriot/kv-storage/internal/storage"
	"github.com/dimishpatriot/kv-storage/internal/storage/localstorage"
	"github.com/gorilla/mux"
)

type App struct {
	logger     *log.Logger
	storage    storage.Storage
	dataLogger transactionlogger.TransactionLogger
	keyService keyservice.KeyService
	handler    handler.Handler
	router     *mux.Router
}

func New() (*App, error) {
	logger := log.New(os.Stdout, "INFO:", log.Lshortfile|log.Ltime|log.Lmicroseconds|log.Ldate)
	logger.Println("logger created")

	// TODO: add selection of storage & logger ->
	storage := localstorage.New()
	logger.Println("storage created")

	dataLogger, err := filelogger.New(logger, "transaction.log", storage)
	if err != nil {
		return nil, fmt.Errorf("failed to create data logger: %w", err)
	}
	logger.Println("dataLogger created")
	// TODO: <------------

	keyService := keyservice.New(logger, storage, dataLogger)
	handler := handler.New(keyService)
	logger.Println("handler created")

	router := mux.NewRouter()
	logger.Println("router created")

	return &App{logger, storage, dataLogger, keyService, handler, router}, nil
}

func (app *App) Run() error {
	app.dataLogger.Run()
	app.logger.Println("dataLogger ran")

	app.addRoutes()
	app.logger.Println("routes added")

	return http.ListenAndServe(":8080", app.router)
}

func (app *App) addRoutes() {
	app.router.HandleFunc("/v1/{key}", app.handler.Put).Methods("PUT")
	app.router.HandleFunc("/v1/{key}", app.handler.Get).Methods("GET")
	app.router.HandleFunc("/v1/{key}", app.handler.Delete).Methods("DELETE")
}
