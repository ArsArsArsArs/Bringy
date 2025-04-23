package database

import "go.mongodb.org/mongo-driver/v2/mongo"

type Database struct {
	DB *mongo.Database
}
