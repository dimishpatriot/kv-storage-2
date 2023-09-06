package main

import (
	"log"

	"github.com/dimishpatriot/kv-storage/cmd/app"
)

func main() {
	app, err := app.New()
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(app.Run())
}
