package commands

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func Start(ctx context.Context, b *bot.Bot, upd *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    upd.Message.Chat.ID,
		ParseMode: "HTML",
		Text:      "Привет! Я простой бот, которого можно добавить в группу, и я буду держать в курсе о последних обсуждениях в группе.\n\n<b>Как это работает?</b>\n1. Вы добавляете меня в группу\n\n2. Настраиваете, устанавливая токен Gemini и добавляя те топики, в которых хотите включить отслеживание\n\n3. Я отправляю сообщение и закрепляю его. Затем я буду обновлять содержание этого сообщения по мере того, как продвигается общение. Участники группы всегда могут взглянуть на закреп и вникнуть в суть дела",
	})
}
