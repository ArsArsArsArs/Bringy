package commands

import (
	"Bringy/services/database"
	"Bringy/services/gemini"
	"Bringy/services/helpful"
	"context"
	"log"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

var tokenSettersMap map[int64]int64

func SetGeminiToken(ctx context.Context, b *bot.Bot, upd *models.Update) {
	if tokenSettersMap == nil {
		tokenSettersMap = make(map[int64]int64)
	}

	if upd.Message.Chat.Type == models.ChatTypeGroup {
		if can := helpful.CanChangeInfo(upd.Message.Chat.ID, upd.Message.From.ID, b); !can {
			return
		}
		tokenSettersMap[upd.Message.From.ID] = upd.Message.Chat.ID

		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: upd.Message.Chat.ID,
			ReplyParameters: &models.ReplyParameters{
				MessageID:                upd.Message.ID,
				AllowSendingWithoutReply: true,
			},
			ParseMode: "HTML",
			Text:      "❓ <b>Зачем нужен токен Gemini?</b>\nМы побуждаем администраторов групп использовать свой токен для того, чтобы ИИ мог обрабатывать сообщения без задержек. Через токен мы используем бесплатную модель Gemini, так что вы можете даже не привязывать банковскую карту\n\n❓ <b>Но где мне его получить?</b>\nВсе подробности здесь: https://ai.google.dev/gemini-api/docs/api-key\n\n❓ <b>Токен у меня, что дальше?</b>\nЯ вас запомнил. Пришлите мне в ЛС эту же команду, а затем туда же - токен. Ни в коем случае не отправляйте токен здесь",
		})
	} else if upd.Message.Chat.Type == models.ChatTypePrivate {
		groupID, ok := tokenSettersMap[upd.Message.From.ID]
		if !ok {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: upd.Message.Chat.ID,
				ReplyParameters: &models.ReplyParameters{
					MessageID:                upd.Message.ID,
					AllowSendingWithoutReply: true,
				},
				ParseMode: "HTML",
				Text:      "Сначала эту команду следует использовать в настраиваемом чате",
			})
			return
		}

		if can := helpful.CanChangeInfo(groupID, upd.Message.From.ID, b); !can {
			return
		}

		splitMsg := strings.Split(upd.Message.Text, " ")
		if len(splitMsg) < 2 {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: upd.Message.Chat.ID,
				ReplyParameters: &models.ReplyParameters{
					MessageID:                upd.Message.ID,
					AllowSendingWithoutReply: true,
				},
				ParseMode: "HTML",
				Text:      "Мы близко! Теперь отправьте не просто эту команду, а команду с токеном. Как-то вот так:\n<code>/set_gemini_token 3fgv23d94v...</code>",
			})
			return
		}

		token := splitMsg[1]

		if ok := gemini.ClientManager.CheckAvailability(token); !ok {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: upd.Message.Chat.ID,
				ReplyParameters: &models.ReplyParameters{
					MessageID:                upd.Message.ID,
					AllowSendingWithoutReply: true,
				},
				ParseMode: "HTML",
				Text:      "Токен не прошёл проверку\n\nИнформация о токенах Gemini: https://ai.google.dev/gemini-api/docs/api-key",
			})
			return
		}

		err := database.DB.SaveGeminiToken(upd.Message.Chat.ID, token)
		if err != nil {
			log.Printf("[ERROR] saving a gemini token. Error: %v", err)
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: upd.Message.Chat.ID,
				ReplyParameters: &models.ReplyParameters{
					MessageID:                upd.Message.ID,
					AllowSendingWithoutReply: true,
				},
				ParseMode: "HTML",
				Text:      "Не удалось сохранить токен. Попробуйте ещё раз, пожалуйста!",
			})
			return
		}

		b.SetMessageReaction(ctx, &bot.SetMessageReactionParams{
			ChatID:    upd.Message.Chat.ID,
			MessageID: upd.Message.ID,
			Reaction: []models.ReactionType{
				{
					Type: models.ReactionTypeTypeEmoji,
					ReactionTypeEmoji: &models.ReactionTypeEmoji{
						Type:  "emoji",
						Emoji: "👍",
					},
				},
			},
		})
	}
}
