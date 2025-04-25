package database

import (
	"Bringy/services/helpful"
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var DB Database

func (db *Database) Connect() {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(fmt.Sprintf("mongodb+srv://%s:%s@cluster0.zdlrkh7.mongodb.net/?retryWrites=true&w=majority&appName=%s", helpful.GetEnvParam("MongoDBUsername", true), helpful.GetEnvParam("MongoDBPassword", true), helpful.GetEnvParam("MongoDBAppName", true))).SetServerAPIOptions(serverAPI).SetTimeout(time.Second * 15)

	client, err := mongo.Connect(opts)
	if err != nil {
		log.Fatalf("[ERROR] connecting to the database. Error: %v", err)
	}

	db.DB = client.Database(helpful.GetEnvParam("MongoDBDatabaseName", true))
	log.Println("[INFO] Connection to the database is established")

	collections, err := db.DB.ListCollectionNames(context.Background(), bson.M{})
	if err != nil {
		log.Fatalf("[ERROR] getting the list of collections. Error: %v", err)
	}

	listOfNecessaryCollections := make(map[string]bool)
	listOfNecessaryCollections["groups"] = false

	for _, collection := range collections {
		listOfNecessaryCollections[collection] = true
	}
	for collectionName, exists := range listOfNecessaryCollections {
		if !exists {
			err := db.DB.CreateCollection(context.Background(), collectionName)
			if err != nil {
				log.Fatalf("[ERROR] creating collection %s. Error: %v", collectionName, err)
			}
			log.Printf("[INFO] Collection %s has been created", collectionName)
		}
	}
}

func (db *Database) SaveGeminiToken(chatID int64, token string) error {
	_, err := db.DB.Collection("groups").UpdateOne(context.Background(), bson.D{{Key: "groupID", Value: chatID}}, bson.D{{Key: "$set", Value: bson.D{{Key: "geminiToken", Value: token}}}}, options.UpdateOne().SetUpsert(true))
	if err != nil {
		return err
	}
	return nil
}

func (db *Database) GetGroupParams(chatID int64) (*Group, error) {
	var group Group
	err := db.DB.Collection("groups").FindOne(context.Background(), bson.D{{Key: "groupID", Value: chatID}}).Decode(&group)
	if err != nil {
		return nil, err
	}
	if group.ID == 0 {
		return nil, errors.New("the group is not found")
	}

	return &group, nil
}

func (db *Database) PutActiveThreadIntoDB(chatID int64, msgID, threadID int, isAlreadyFound bool) error {
	if isAlreadyFound {
		_, err := db.DB.Collection("groups").UpdateOne(
			context.Background(),
			bson.D{{Key: "groupID", Value: chatID}},
			bson.D{{Key: "$set", Value: bson.M{
				"threads.$[thread].active":          true,
				"threads.$[thread].pinnedMessageID": msgID,
			}}},
			options.UpdateOne().SetUpsert(true).SetArrayFilters([]interface{}{
				bson.M{"thread.threadID": threadID},
			}),
		)
		if err != nil {
			return err
		}
	} else {
		_, err := db.DB.Collection("groups").UpdateOne(context.Background(), bson.D{{Key: "groupID", Value: chatID}}, bson.D{{Key: "$push", Value: bson.D{{Key: "threads", Value: GroupThread{
			ThreadID:        threadID,
			Active:          true,
			PinnedMessageID: msgID,
		}}}}}, options.UpdateOne().SetUpsert(true))
		if err != nil {
			return err
		}
	}
	return nil
}

func (db *Database) PullActiveThreadOutOfDB(chatID int64, threadID int) error {
	_, err := db.DB.Collection("groups").UpdateOne(
		context.Background(),
		bson.D{{Key: "groupID", Value: chatID}},
		bson.D{{Key: "$set", Value: bson.M{
			"threads.$[thread].active": false,
		}}},
		options.UpdateOne().SetArrayFilters([]interface{}{
			bson.M{"thread.threadID": threadID},
		}),
	)
	if err != nil {
		return err
	}
	return nil
}
