package commands

import (
	"Bringy/services/database"
	"Bringy/services/helpful"
	"context"
	"fmt"
	"log"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func Launch(ctx context.Context, b *bot.Bot, upd *models.Update) {
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

	var found bool
	var foundThread database.GroupThread
	for _, thread := range group.Threads {
		if thread.ThreadID == upd.Message.MessageThreadID {
			found = true
			foundThread = thread
			break
		}
	}

	if foundThread.Active {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: upd.Message.Chat.ID,
			ReplyParameters: &models.ReplyParameters{
				MessageID:                upd.Message.ID,
				AllowSendingWithoutReply: true,
			},
			ParseMode: "HTML",
			Text:      fmt.Sprintf("Я уже работаю в этом месте\n\n👀 <a href=\"https://t.me/c/%d/%d\">Посмотреть сводку</a>\n\nЧто-то идёт не так? Используйте команду /stop, а затем повторно /launch", helpful.InternalGroupID(upd.Message.Chat.ID), foundThread.PinnedMessageID),
		})
		return
	}

	msg, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:          upd.Message.Chat.ID,
		MessageThreadID: upd.Message.MessageThreadID,
		ParseMode:       "HTML",
		Text:            "<i>Это сообщение будет использоваться для отображения сводки</i>",
	})
	if err != nil {
		log.Printf("[ERROR] sending a message for summarizations. Error: %v", err)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: upd.Message.Chat.ID,
			ReplyParameters: &models.ReplyParameters{
				MessageID:                upd.Message.ID,
				AllowSendingWithoutReply: true,
			},
			ParseMode: "HTML",
			Text:      "При отправке служебного сообщения произошла ошибка",
		})
		return
	}
	_, err = b.PinChatMessage(ctx, &bot.PinChatMessageParams{
		ChatID:              upd.Message.Chat.ID,
		MessageID:           msg.ID,
		DisableNotification: true,
	})
	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: upd.Message.Chat.ID,
			ReplyParameters: &models.ReplyParameters{
				MessageID:                upd.Message.ID,
				AllowSendingWithoutReply: true,
			},
			ParseMode: "HTML",
			Text:      "Служебное сообщение отправлено, но его не удалось закрепить. Рекомендуется сделать это самостоятельно",
		})
	}

	err = database.DB.PutActiveThreadIntoDB(upd.Message.Chat.ID, msg.ID, upd.Message.MessageThreadID, found)
	if err != nil {
		log.Printf("[ERROR] putting an active thread into DB. Error: %v", err)
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

	if group.GeminiToken == "" {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: upd.Message.Chat.ID,
			ReplyParameters: &models.ReplyParameters{
				MessageID:                upd.Message.ID,
				AllowSendingWithoutReply: true,
			},
			ParseMode: "HTML",
			Text:      "Токен Gemini ещё не был установлен для этой группы. Без него я не смогу работать. Подробнее: /set_gemini_token",
		})
	}
}
