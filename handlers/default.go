package handlers

import (
	"context"
	"log"
	"strconv"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func DefaultHandler(ctx context.Context, b *bot.Bot, upd *models.Update) {
	joinedOrLeftGroup(upd)
}

func joinedOrLeftGroup(upd *models.Update) {
	if upd.MyChatMember != nil {
		idStr := strconv.Itoa(int(upd.MyChatMember.Chat.ID))
		if (upd.MyChatMember.NewChatMember.Banned != nil || upd.MyChatMember.NewChatMember.Left != nil) && (upd.MyChatMember.OldChatMember.Member != nil || upd.MyChatMember.OldChatMember.Administrator != nil || upd.MyChatMember.OldChatMember.Owner != nil) && strings.HasPrefix(idStr, "-") {
			log.Printf("[INFO] A group \"%s\" (%d) was removed", upd.MyChatMember.Chat.Title, upd.MyChatMember.Chat.ID)
		}
		if (upd.MyChatMember.OldChatMember.Banned != nil || upd.MyChatMember.OldChatMember.Left != nil) && (upd.MyChatMember.NewChatMember.Member != nil || upd.MyChatMember.NewChatMember.Administrator != nil || upd.MyChatMember.NewChatMember.Owner != nil) && strings.HasPrefix(idStr, "-") {
			log.Printf("[INFO] A group \"%s\" (%d) was added", upd.MyChatMember.Chat.Title, upd.MyChatMember.Chat.ID)
		}
	}
}
