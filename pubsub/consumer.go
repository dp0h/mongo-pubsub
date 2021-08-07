package pubsub

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type Handler func(ctx context.Context, opts ...interface{}) error

func (h Handler) Handle(ctx context.Context, opts ...interface{}) error {
	return h(ctx, opts...)
}

type Consumer interface {
	Handle(ctx context.Context, opts ...interface{}) error
	Topic() string
	DB() *mongo.Database
}

type consumer struct {
	topic   string
	handler Handler
	db      *mongo.Database
}

func (c consumer) Handle(ctx context.Context, opts ...interface{}) error {
	err := c.handler(ctx, opts...)
	return err
}

func (c consumer) Topic() string {
	return c.topic
}

func (c consumer) DB() *mongo.Database {
	return c.db
}

func NewConsumer(topic string, handler Handler, db *mongo.Database) Consumer {
	return &consumer{
		topic:   topic,
		handler: handler,
		db:      db,
	}
}
