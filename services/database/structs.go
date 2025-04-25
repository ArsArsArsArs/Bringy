package database

import "go.mongodb.org/mongo-driver/v2/mongo"

type Database struct {
	DB *mongo.Database
}

type Group struct {
	ID          int64         `bson:"groupID"`
	Threads     []GroupThread `bson:"threads"`
	GeminiToken string        `bson:"geminiToken"`
}
type GroupThread struct {
	ThreadID        int  `bson:"threadID"`
	Active          bool `bson:"active"`
	PinnedMessageID int  `bson:"pinnedMessageID"`
}
