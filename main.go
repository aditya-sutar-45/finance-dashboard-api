package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/aditya-sutar-45/finance-dashboard-api/app"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("error loading .env ", err)
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatalln("DB_URL is not found in the env")
	}
	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatalln("PORT is not found in the env")
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	app, err := app.New(portString, dbURL)
	if err != nil {
		log.Println("failed to create the app", err)
	}

	err = app.Start(ctx)
	if err != nil {
		log.Println(err)
	}
}
