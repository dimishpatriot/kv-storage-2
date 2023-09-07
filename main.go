package main

import (
	"log"
	"os"

	"github.com/dimishpatriot/kv-storage/cmd/app"
)

func main() {
	app, err := app.New()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	log.Fatal(app.Run())
}
