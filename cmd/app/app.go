package app

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/dimishpatriot/kv-storage/internal/handler"
	"github.com/dimishpatriot/kv-storage/internal/services/keyservice"
	"github.com/dimishpatriot/kv-storage/internal/services/transactionlogger"
	"github.com/dimishpatriot/kv-storage/internal/services/transactionlogger/filelogger"
	"github.com/dimishpatriot/kv-storage/internal/services/transactionlogger/postgreslogger"
	"github.com/dimishpatriot/kv-storage/internal/storage"
	"github.com/dimishpatriot/kv-storage/internal/storage/localstorage"
	"github.com/dimishpatriot/kv-storage/internal/storage/postgresstorage"
	"github.com/gorilla/mux"
)

type App struct {
	logger     *log.Logger
	dataLogger transactionlogger.TransactionLogger
	keyService keyservice.KeyService
	handler    handler.Handler
	storage    storage.Storage
	router     *mux.Router
}

type AppConfig struct {
	StorageType string
}

var (
	LocalStorage = "local"
	PGStorage    = "postgres"
)

func New(config AppConfig) (*App, error) {
	var storage storage.Storage
	var dataLogger transactionlogger.TransactionLogger
	var err error
	var db *sql.DB

	logger := log.New(os.Stdout, "INFO:", log.Lshortfile|log.Ltime|log.Lmicroseconds|log.Ldate)
	logger.Println("logger created")

	switch config.StorageType {

	case LocalStorage:
		storage = localstorage.New()
		logger.Println("storage created")

		dataLogger, err = filelogger.New(logger, "transaction.log")
		if err != nil {
			return nil, fmt.Errorf("failed to create file-logger: %w", err)
		}
		logger.Println("dataLogger created")
		restoreData(dataLogger, storage)
		logger.Println("data restored")

	case PGStorage:
		dbParams := postgreslogger.PostgresDBParams{
			Host:     os.Getenv("DB_HOST"),
			DBName:   os.Getenv("DB_NAME"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			SSLMode:  os.Getenv("DB_SSL_MODE"),
		}

		dataLogger, db, err = postgreslogger.New(logger, dbParams)
		if err != nil {
			return nil, fmt.Errorf("failed to create pg-logger: %w", err)
		}
		storage = postgresstorage.New(db)
		logger.Println("storage created")

		logger.Println("dataLogger created")

	default:
		return nil, fmt.Errorf("invalid type of storage: %s", config.StorageType)
	}

	keyService := keyservice.New(logger, storage, dataLogger)
	logger.Println("keyservice created")

	handler := handler.New(keyService)
	logger.Println("handler created")

	router := mux.NewRouter()
	logger.Println("router created")

	return &App{logger, dataLogger, keyService, handler, storage, router}, nil
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

func restoreData(
	fileLogger transactionlogger.TransactionLogger,
	storage storage.Storage,
) {
	var err error
	events, errors := fileLogger.ReadEvents()
	e, ok := transactionlogger.Event{}, true

	for ok && err == nil {
		select {
		case err, ok = <-errors:
		case e, ok = <-events:
			switch e.EventType {
			case transactionlogger.EventDelete:
				err = storage.Delete(e.Key)
			case transactionlogger.EventPut:
				err = storage.Put(e.Key, e.Value)
			}
		}
	}
}
