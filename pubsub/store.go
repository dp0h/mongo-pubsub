package pubsub

import (
	"context"
	"log"
	"strings"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Global TODO: global variables cause package blocks
// find a different way
var Global struct {
	CounterMap map[string]int64
	CounterCh  chan string
	InitCh     chan map[string]int64
}

type PubSub struct {
	Client     *mongo.Client
	AppDb      *string
	publishers map[string]Publisher
	consumers  map[string]Consumer
	lock       *sync.RWMutex
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
		Client:     client,
		AppDb:      &appDb,
		lock:       &sync.RWMutex{},
		publishers: make(map[string]Publisher),
		consumers:  make(map[string]Consumer),
	}

	return &session, nil
}

func (p *PubSub) Collection(name string) *mongo.Collection {
	return p.Client.Database(*p.AppDb).Collection(name)
}

func (p *PubSub) RegisterPublishers(topics ...string) {
	var (
		capped       = true
		size   int64 = 100000
		err    error
	)

	p.setPublishers(topics...)

	for k := range p.publishers {
		err = p.Client.Database(*p.AppDb).CreateCollection(context.TODO(), k, &options.CreateCollectionOptions{
			Capped:      &capped,
			SizeInBytes: &size,
		})

		if err == nil {
			continue
		}
		if err != nil && !strings.Contains(err.Error(), "NamespaceExists") {
			log.Fatalf("error occurred: %v\n", err)
		}

		myOptions := options.FindOne()
		myOptions.SetSort(bson.M{"$natural": -1})
		collection := p.Client.Database(*p.AppDb).Collection(k)
		result := collection.FindOne(context.TODO(), bson.M{}, myOptions)
		event := &Event{}
		err = result.Decode(&event)

		m := map[string]int64{
			k: event.ID,
		}
		Global.InitCh <- m

	}
}

func (p *PubSub) setPublishers(topics ...string) {
	p.lock.Lock()
	defer p.lock.Unlock()

	for _, topic := range topics {
		p.publishers[topic] = NewPublisher(topic, p.Client.Database(*p.AppDb))
	}

}

func (p *PubSub) AddPublisher(topic string) bool {
	p.lock.Lock()
	defer p.lock.Unlock()

	var ok bool
	if _, ok = p.publishers[topic]; !ok {
		p.publishers[topic] = NewPublisher(topic, p.Client.Database(*p.AppDb))
	}
	return ok
}

func (p *PubSub) Publishers() map[string]Publisher {
	return p.publishers
}

func (p *PubSub) Consumers() map[string]Consumer {
	return p.consumers
}

func (p *PubSub) AddConsumer(topic string, handler Handler, db *mongo.Database) bool {
	consumer := NewConsumer(topic, handler, db)

	p.lock.RLock()
	_, ok := p.consumers[topic]
	p.lock.RUnlock()

	if !ok {
		p.lock.Lock()
		p.consumers[topic] = consumer
		p.lock.Unlock()
	}

	return ok
}
