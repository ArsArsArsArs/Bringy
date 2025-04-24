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
	"github.com/go-telegram/bot/models"
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
	go registerCommands(ctx, b)
	b.Start(ctx)
}

func registerCommands(ctx context.Context, b *bot.Bot) {
	_, err := b.SetMyCommands(ctx, &bot.SetMyCommandsParams{
		Commands: []models.BotCommand{
			{
				Command:     "start",
				Description: "Базовая информация",
			},
			{
				Command:     "set_gemini_token",
				Description: "Сначала эту команду нужно ввести в группе",
			},
		},
		Scope: &models.BotCommandScopeAllPrivateChats{},
	})
	if err != nil {
		log.Printf("[WARNING] Private chats' commands haven't been added. Reason: %v", err)
		err = nil
	}

	_, err = b.SetMyCommands(ctx, &bot.SetMyCommandsParams{
		Commands: []models.BotCommand{
			{
				Command:     "launch",
				Description: "Запускает работу в текущем топике",
			},
			{
				Command:     "stop",
				Description: "Останавливает работу в текущем топике",
			},
			{
				Command:     "set_gemini_token",
				Description: "❗ Установить токен ИИ. Требуется для работы",
			},
		},
		Scope: &models.BotCommandScopeChatAdministrators{},
	})
	if err != nil {
		log.Printf("[WARNING] Commands for chat administators haven't been added. Reason: %v", err)
		err = nil
	}
}
