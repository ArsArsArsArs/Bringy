package main

import (
	"Bringy/handlers"
	"Bringy/services/database"
	"Bringy/services/helpful"
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/go-telegram/bot"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("[ERROR] loading godotenv. Error: %v", err)
	}

	database.DB.Connect()
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithErrorsHandler(handlers.ErrorHandler),
		bot.WithDefaultHandler(handlers.DefaultHandler),
	}

	b, err := bot.New(helpful.GetEnvParam("BotToken", true), opts...)
	if err != nil {
		log.Fatalf("[ERROR] initiating a new bot. Error: %v", err)
	}

	log.Println("[INFO] Bringy should be started")
	b.Start(ctx)
}
