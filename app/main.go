package main

import (
	"github.com/dp0h/mongo-pubsub/pubsub"
	"github.com/dp0h/mongo-pubsub/rest"
	"github.com/jessevdk/go-flags"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"time"
)

var opts struct {
	Host       string `long:"host" env:"APP_HOST" required:"true" description:"host"`
	Port       int    `long:"port" env:"APP_PORT" required:"true" description:"port"`
	MongoDbURI string `long:"mongodb-uri" env:"APP_MONGODB_URI" required:"true" description:"mongodb uri"`
	AppDb      string `long:"appdb" env:"APP_DB" required:"true" description:"mongodb app db"`
}

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	zerolog.SetGlobalLevel(zerolog.TraceLevel)

	if _, err := flags.Parse(&opts); err != nil {
		log.Fatal().Err(err).Msg("failed to parse args")
	}

	pubSub, err := pubsub.NewPubSub(opts.MongoDbURI, opts.AppDb)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to mongodb")
	}

	srv, err := rest.NewRestSrv(pubSub)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create http server")
	}

	err = srv.ListenAndServe(opts.Host, opts.Port)
	if err != nil {
		log.Error().Err(err).Msg("Failed to ListenAndServe")
	}
}
