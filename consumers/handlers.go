package consumers

import (
	"context"
	"log"

	pb "github.com/dp0h/mongo-pubsub/pubsub"
)

// triggered for every instance
func postEventDoSomething(ctx context.Context, opts ...interface{}) error {
	log.Printf(
		"Well done, it is working!\n ID: %v, Message: %v, Time: %v\n",
		opts[0].(pb.Event).ID,
		opts[0].(pb.Event).Message,
		opts[0].(pb.Event).Time,
	)
	return nil
}

func getEventsDoSomething(ctx context.Context, opts ...interface{}) error {
	// Do something
	return nil
}
