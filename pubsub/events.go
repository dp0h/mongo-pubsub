package pubsub

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Event struct {
	Message string
	Time     time.Time
}

const (
	eventsTable     string = "events"
)

func (p *PubSub) AddEvent(message string) error {
	newEvent := Event{
		Message: message,
		Time:     time.Now().UTC(),
	}
	_, err := p.Collection(eventsTable).InsertOne(context.TODO(), newEvent)
	return err
}

func (p *PubSub) GetEvents() ([]Event, error) {
	cur, err := p.Collection(eventsTable).Find(context.TODO(), bson.D{{}}, options.Find())
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
