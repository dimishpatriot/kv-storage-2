package app

import (
	"fmt"
	"net/http"

	"github.com/dimishpatriot/kv-storage/internal/controller"
	"github.com/dimishpatriot/kv-storage/internal/logger/datalogger"
	"github.com/gorilla/mux"
)

type App struct {
	dataLogger datalogger.TransactionLogger
	router     *mux.Router
	controller *controller.Controller
}

func New() (*App, error) {
	dl, err := datalogger.New("transaction.log")
	if err != nil {
		return nil, fmt.Errorf("failed to create data logger: %w", err)
	}

	c := controller.New(dl)
	r := mux.NewRouter()
	addRoutes(r, c)

	return &App{dataLogger: dl, router: r, controller: c}, nil
}

func addRoutes(r *mux.Router, c *controller.Controller) {
	r.HandleFunc("/v1/{key}", c.PutHandler).Methods("PUT")
	r.HandleFunc("/v1/{key}", c.GetHandler).Methods("GET")
	r.HandleFunc("/v1/{key}", c.DeleteHandler).Methods("DELETE")
}

func (a *App) Run() error {
	a.dataLogger.RestoreData(a.controller.Storage)
	a.dataLogger.Run()

	return http.ListenAndServe(":8080", a.router)
}
