package pubsub

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type PubSub struct {
	Client *mongo.Client
	AppDb  *string
}

func NewPubSub(mongoDbURI string, appDb string) (*PubSub, error) {
	clientOptions := options.Client().ApplyURI(mongoDbURI).SetReadPreference(readpref.Nearest())

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, err
	}

	session := PubSub{
		Client: client,
		AppDb:  &appDb,
	}

	return &session, nil
}

func (p *PubSub) Collection(name string) *mongo.Collection {
	return p.Client.Database(*p.AppDb).Collection(name)
}
