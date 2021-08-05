package main

import (
	"github.com/jessevdk/go-flags"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"time"
)

var opts struct {
	MongoDbURI string `long:"mongodb-uri" env:"APP_MONGODB_URI" required:"true" description:"mongodb uri"`
	AppDb      string `long:"appdb" env:"MONGO_APP_DB" required:"true" description:"mongodb app db"`
}

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	zerolog.SetGlobalLevel(zerolog.TraceLevel)

	if _, err := flags.Parse(&opts); err != nil {
		log.Fatal().Err(err).Msg("failed to parse args")
	}

}
