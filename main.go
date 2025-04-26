package main

import (
	"Bringy/handlers"
	"Bringy/handlers/commands"
	"Bringy/services/database"
	"Bringy/services/gemini"
	"Bringy/services/helpful"
	"Bringy/services/summarization"
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
		log.Printf("[WARNING] loading godotenv. Reason: %v", err)
	}

	database.DB.Connect()
	gemini.NewClientManager()
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithErrorsHandler(handlers.ErrorHandler),
		bot.WithDefaultHandler(handlers.DefaultHandler),
		bot.WithAllowedUpdates(bot.AllowedUpdates{"message", "my_chat_member"}),
	}

	b, err := bot.New(helpful.GetEnvParam("BotToken", true), opts...)
	if err != nil {
		log.Fatalf("[ERROR] initiating a new bot. Error: %v", err)
	}

	registerHandlers(b)

	log.Println("[INFO] Bringy should be started")
	go registerCommands(ctx, b)
	go summarization.StartingInit(b)
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
		Scope: &models.BotCommandScopeAllChatAdministrators{},
	})
	if err != nil {
		log.Printf("[WARNING] Commands for chat administators haven't been added. Reason: %v", err)
		err = nil
	}
}

func registerHandlers(b *bot.Bot) {
	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypePrefix, commands.Start, privateOnly)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/set_gemini_token", bot.MatchTypePrefix, commands.SetGeminiToken)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/launch", bot.MatchTypePrefix, commands.Launch, groupsOnly, adminsOnly)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/stop", bot.MatchTypePrefix, commands.Stop, groupsOnly, adminsOnly)
}

func privateOnly(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, upd *models.Update) {
		if upd.Message.Chat.Type != models.ChatTypePrivate {
			return
		}
		next(ctx, b, upd)
	}
}

func groupsOnly(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, upd *models.Update) {
		if (upd.Message.Chat.Type != models.ChatTypeGroup) && (upd.Message.Chat.Type != models.ChatTypeSupergroup) {
			return
		}
		next(ctx, b, upd)
	}
}

func adminsOnly(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, upd *models.Update) {
		if can := helpful.CanChangeInfo(upd.Message.Chat.ID, upd.Message.From.ID, b); !can {
			return
		}
		next(ctx, b, upd)
	}
}
