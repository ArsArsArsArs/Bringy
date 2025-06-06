package summarization

import (
	"Bringy/services/config"
	"Bringy/services/database"
	"Bringy/services/gemini"
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/go-telegram/bot"
)

func summarize(cb *CircularBuffer) {
	msgs := cb.GetAll()

	if cb.count == 0 {
		return
	}

	group, err := database.DB.GetGroupParams(cb.groupID)
	if err != nil {
		return
	}

	if group.GeminiToken == "" {
		return
	}

	toSummarize := strings.Join(msgs, "[NEXT MESSAGE]")
	summarization, err := gemini.SummarizeMessages(group.GeminiToken, toSummarize)
	if err != nil {
		log.Printf("[ERROR] summarizing. Error: %v", err)
		return
	}

	timeNow := time.Now().UTC().Add(time.Hour * time.Duration(config.UTCPlusHours))

	b.EditMessageText(context.Background(), &bot.EditMessageTextParams{
		ChatID:    group.ID,
		MessageID: cb.pinnedMessageID,
		ParseMode: "HTML",
		Text:      fmt.Sprintf("👀 <b>О чём идёт речь сейчас?</b>\n\n<blockquote>%s</blockquote>\n<i>Последнее обновление: %s (UTC+%d)</i>", summarization, timeNow.Format(time.TimeOnly), config.UTCPlusHours),
	})

	cb.Clear()
}
