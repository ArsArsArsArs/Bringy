package summarization

import (
	"Bringy/services/database"
	"log"
	"strconv"
	"time"

	"github.com/go-telegram/bot"
)

var Buffers map[string]*CircularBuffer
var b *bot.Bot

func AddBuffer(chatID int64, threadID int, pinnedMessageID int) *CircularBuffer {
	if Buffers == nil {
		Buffers = make(map[string]*CircularBuffer)
	}

	chatIDstr := strconv.Itoa(int(chatID))
	threadIDstr := strconv.Itoa(threadID)

	newCb := NewCircularBuffer(60, chatID, threadID, pinnedMessageID)

	Buffers[chatIDstr+"_"+threadIDstr] = newCb
	return newCb
}

func ProcessBuffer(cb *CircularBuffer) {
	ticker := time.NewTicker(time.Minute * 5)
	defer ticker.Stop()

	for range ticker.C {
		summarize(cb)
	}
}

func StartingInit(bot *bot.Bot) {
	b = bot
	log.Println("[INFO] Starting buffers initialization")
	groups, err := database.DB.GetActiveGroups()
	if err != nil {
		log.Fatalf("[ERROR] getting active groups. Error: %v", err)
	}
	log.Printf("[INFO] Active groups fetched: %d", len(*groups))

	var activeThreads int
	for _, group := range *groups {
		for _, thread := range group.Threads {
			if thread.Active {
				AddBuffer(group.ID, thread.ThreadID, thread.PinnedMessageID)
				activeThreads++
			}
		}
	}

	log.Printf("[INFO] Active threads fetched: %d", activeThreads)

	for _, bufferItself := range Buffers {
		go ProcessBuffer(bufferItself)
	}
}
