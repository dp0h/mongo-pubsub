package pubsub

import "go.mongodb.org/mongo-driver/mongo"

type PubSub struct {
	MongoClient *mongo.Client
}

