package commands

import (
	"Bringy/services/database"
	"context"
	"log"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func Stop(ctx context.Context, b *bot.Bot, upd *models.Update) {
	group, err := database.DB.GetGroupParams(upd.Message.Chat.ID)
	if err != nil {
		log.Printf("[ERROR] looking for a group in launch. Group ID: %d. Error: %v", upd.Message.Chat.ID, err)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: upd.Message.Chat.ID,
			ReplyParameters: &models.ReplyParameters{
				MessageID:                upd.Message.ID,
				AllowSendingWithoutReply: true,
			},
			ParseMode: "HTML",
			Text:      "При поиске этой группы в базе данных произошла ошибка",
		})
		return
	}

	var foundThread database.GroupThread
	for _, thread := range group.Threads {
		if thread.ThreadID == upd.Message.MessageThreadID {
			foundThread = thread
			break
		}
	}

	if !foundThread.Active {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: upd.Message.Chat.ID,
			ReplyParameters: &models.ReplyParameters{
				MessageID:                upd.Message.ID,
				AllowSendingWithoutReply: true,
			},
			ParseMode: "HTML",
			Text:      "Здесь и так не работает сводка. Её можно запустить с помощью команды /launch",
		})
		return
	}

	err = database.DB.PullActiveThreadOutOfDB(upd.Message.Chat.ID, upd.Message.MessageThreadID)
	if err != nil {
		log.Printf("[ERROR] pulling an active thread out of DB. Error: %v", err)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: upd.Message.Chat.ID,
			ReplyParameters: &models.ReplyParameters{
				MessageID:                upd.Message.ID,
				AllowSendingWithoutReply: true,
			},
			ParseMode: "HTML",
			Text:      "При сохранении информации в базу данных произошла ошибка. Рекомендуется ввести команду ещё раз",
		})
		return
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: upd.Message.Chat.ID,
		ReplyParameters: &models.ReplyParameters{
			MessageID:                upd.Message.ID,
			AllowSendingWithoutReply: true,
		},
		ParseMode: "HTML",
		Text:      "✅",
	})
}
