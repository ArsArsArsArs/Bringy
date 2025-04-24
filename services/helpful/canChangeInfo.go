package helpful

import (
	"context"
	"log"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func CanChangeInfo(chatID int64, userID int64, b *bot.Bot) bool {
	listOfAdministrators, err := b.GetChatAdministrators(context.Background(), &bot.GetChatAdministratorsParams{
		ChatID: chatID,
	})
	if err != nil {
		log.Printf("[WARNING] getting the list of administators of the group with ID %d. Reason: %v", chatID, err)
		return false
	}
	for _, administrator := range listOfAdministrators {
		switch administrator.Type {
		case models.ChatMemberTypeOwner:
			if administrator.Owner.User.ID == userID {
				return true
			}
		case models.ChatMemberTypeAdministrator:
			if administrator.Administrator.User.ID == userID {
				if administrator.Administrator.CanChangeInfo {
					return true
				} else {
					return false
				}
			}
		}
	}
	return false
}
