package main

import (
	"os"
	"time"

	"github.com/dp0h/mongo-pubsub/consumers"
	"github.com/dp0h/mongo-pubsub/pubsub"
	"github.com/dp0h/mongo-pubsub/rest"
	"github.com/jessevdk/go-flags"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var opts struct {
	Host       string `long:"host" env:"APP_HOST" required:"true" description:"host"`
	Port       int    `long:"port" env:"APP_PORT" required:"true" description:"port"`
	MongoDbURI string `long:"mongodb-uri" env:"APP_MONGODB_URI" required:"true" description:"mongodb uri"`
	AppDb      string `long:"appdb" env:"APP_DB" required:"true" description:"mongodb app db"`
}

func init() {
	pubsub.Global = struct {
		CounterMap map[string]int64
		CounterCh  chan string
		InitCh     chan map[string]int64
	}{
		CounterMap: make(map[string]int64),
		CounterCh:  make(chan string),
		InitCh:     make(chan map[string]int64),
	}

	go func() {
		for {
			select {
			case v := <-pubsub.Global.CounterCh:
				pubsub.Global.CounterMap[v]++
			case m := <-pubsub.Global.InitCh:
				for k, v := range m {
					pubsub.Global.CounterMap[k] = v
				}
			}
		}
	}()
}

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	zerolog.SetGlobalLevel(zerolog.TraceLevel)

	if _, err := flags.Parse(&opts); err != nil {
		log.Fatal().Err(err).
			Msg("failed to parse args")
	}

	pubSub, err := pubsub.NewPubSub(opts.MongoDbURI, opts.AppDb)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to mongodb")
	}

	topics := []string{"v1_mongo_pubsub_event_post", "v1_mongo_pubsub_events_get"}
	pubSub.RegisterPublishers(topics...)

	consumers.Register(*pubSub)
	go consumers.Run(*pubSub)

	srv, err := rest.NewRestSrv(pubSub)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create http server")
	}

	err = srv.ListenAndServe(opts.Host, opts.Port)
	if err != nil {
		log.Error().Err(err).Msg("Failed to ListenAndServe")
	}
}
