package pubsub

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Event struct {
	ID      int64
	Message string
	Time    time.Time
}

const (
	eventsTable string = "events"
)

func (p *PubSub) AddEvent(message string) error {
	newEvent := Event{
		Message: message,
		Time:    time.Now().UTC(),
	}
	_, err := p.Collection(eventsTable).InsertOne(context.TODO(), newEvent)
	return err
}

func (p *PubSub) GetEvents(topic string) ([]Event, error) {
	cur, err := p.Collection(topic).Find(context.TODO(), bson.D{{}}, options.Find())
	if err != nil {
		return nil, err
	}

	results := make([]Event, 0)

	for cur.Next(context.TODO()) {
		var evt Event
		err := cur.Decode(&evt)
		if err != nil {
			return nil, err
		}

		results = append(results, evt)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	_ = cur.Close(context.TODO())

	return results, nil
}
