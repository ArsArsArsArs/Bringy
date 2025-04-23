package database

import (
	"Bringy/services/helpful"
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var DB Database

func (db Database) Connect() {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(fmt.Sprintf("mongodb+srv://%s:%s@cluster0.zdlrkh7.mongodb.net/?retryWrites=true&w=majority&appName=%s", helpful.GetEnvParam("MongoDBUsername", true), helpful.GetEnvParam("MongoDBPassword", true), helpful.GetEnvParam("MongoDBAppName", true))).SetServerAPIOptions(serverAPI).SetTimeout(time.Second * 15)

	client, err := mongo.Connect(opts)
	if err != nil {
		log.Fatalf("[ERROR] connecting to the database. Error: %v", err)
	}

	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			log.Fatalf("[ERROR] disconnecting from the database. Error: %v", err)
		}
	}()

	db.DB = client.Database(helpful.GetEnvParam("MongoDBDatabaseName", true))
	log.Println("[INFO] Connection to the database is established")

	collections, err := db.DB.ListCollectionNames(context.Background(), bson.M{})
	if err != nil {
		log.Fatalf("[ERROR] getting the list of collections. Error: %v", err)
	}

	listOfNecessaryCollections := make(map[string]bool, 0)
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
