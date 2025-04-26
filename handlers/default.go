package handlers

import (
	"Bringy/services/summarization"
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func DefaultHandler(ctx context.Context, b *bot.Bot, upd *models.Update) {
	joinedOrLeftGroup(upd)
	watchingMessages(upd)
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

func watchingMessages(upd *models.Update) {
	if (upd.Message != nil) && (upd.Message.Text != "") && ((upd.Message.Chat.Type == models.ChatTypeGroup) || (upd.Message.Chat.Type == models.ChatTypeSupergroup)) {
		cb, found := summarization.Buffers[fmt.Sprintf("%d_%d", upd.Message.Chat.ID, upd.Message.MessageThreadID)]
		if !found {
			return
		}

		cb.Add(fmt.Sprintf("%s: %s", upd.Message.From.FirstName, upd.Message.Text))
	}
}
