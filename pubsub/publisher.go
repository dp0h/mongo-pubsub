package pubsub

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

func NewPublisher(topic string, db *mongo.Database) Publisher {
	return publisher{
		topic: topic,
		db:    db,
	}
}

type Publisher interface {
	Push(ctx context.Context, data interface{}) error
	Topic() string
}

type publisher struct {
	topic string
	db    *mongo.Database
}

func (p publisher) Push(ctx context.Context, data interface{}) error {
	_, err := p.db.Collection(p.topic).InsertOne(ctx, data)
	return err
}

func (p publisher) Topic() string {
	return p.topic
}
