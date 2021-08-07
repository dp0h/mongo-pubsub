package consumers

import (
	"context"
	"log"
	"time"

	"github.com/dp0h/mongo-pubsub/pubsub"
	pb "github.com/dp0h/mongo-pubsub/pubsub"
	"go.mongodb.org/mongo-driver/bson"
)

func Register(pubsub pb.PubSub) {
	v1Route := "v1_mongo_pubsub"

	pubsub.AddConsumer(
		v1Route+"_event_post",                 // topic name
		postEventDoSomething,                  // handler
		pubsub.Client.Database(*pubsub.AppDb), // db
	)

	pubsub.AddConsumer(
		v1Route+"_events_get",                 // topic name
		postEventDoSomething,                  // handler
		pubsub.Client.Database(*pubsub.AppDb), // db
	)
}

func Run(pubsub pubsub.PubSub) {
	for topic, consumer := range pubsub.Consumers() {
		for {
			ctx := context.TODO()
			cursor, err := consumer.DB().Collection(topic).Find(ctx, bson.D{{"id", bson.D{{"$gt", pb.Global.CounterMap[topic]}}}})
			defer cursor.Close(ctx)
			if err != nil {
				log.Println("error occurred: ", err)
			} else {
				for cursor.Next(ctx) {
					event := pb.Event{}

					if err = cursor.Decode(&event); err != nil {
						return
					}

					consumer.Handle(ctx, event)
					pb.Global.CounterCh <- topic
				}
			}
			<-time.After(time.Millisecond * 100)
		}

	}
}
