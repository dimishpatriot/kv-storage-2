package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/dimishpatriot/kv-storage/internal/controller"
	"github.com/dimishpatriot/kv-storage/internal/logger/datalogger"
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
	dataLogger datalogger.TransactionLogger
	router     *mux.Router
	controller *controller.Controller
}

func New() (*App, error) {
	l := log.New(os.Stdout, "INFO:", log.Lshortfile|log.Ltime|log.Lmicroseconds|log.Ldate)

	dl, err := datalogger.New(l, "transaction.log")
	if err != nil {
		return nil, fmt.Errorf("failed to create data logger: %w", err)
	}
	l.Println("dataLogger created")

	c := controller.New(l, dl)
	l.Println("controller created")

	r := mux.NewRouter()
	l.Println("router created")

	return &App{dataLogger: dl, router: r, controller: c, logger: l}, nil
}

func (a *App) addRoutes() {
	a.router.HandleFunc("/v1/{key}", a.controller.PutHandler).Methods("PUT")
	a.router.HandleFunc("/v1/{key}", a.controller.GetHandler).Methods("GET")
	a.router.HandleFunc("/v1/{key}", a.controller.DeleteHandler).Methods("DELETE")
}

func (a *App) Run() error {
	a.addRoutes()
	a.logger.Println("routes added")

	a.dataLogger.RestoreData(a.controller.Storage)
	a.logger.Println("data restored")

	a.dataLogger.Run()
	a.logger.Println("dataLogger ran")

	return http.ListenAndServe(":8080", a.router)
}
