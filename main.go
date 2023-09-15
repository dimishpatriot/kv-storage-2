package main

import (
	"flag"
	"log"

	"github.com/dimishpatriot/kv-storage/cmd/app"
	"github.com/joho/godotenv"
)

func main() {
	storageType := flag.String("s", "local", "type of storage")
	flag.Parse()

	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("can't get environment variables: %w", err)
	}

	app, err := app.New(app.AppConfig{StorageType: *storageType})
	if err != nil {
		log.Fatal("can't create new application: %w", err)
	}
	log.Fatal(app.Run())
}
