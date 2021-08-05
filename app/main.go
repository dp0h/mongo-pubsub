package main

import (
	"github.com/jessevdk/go-flags"
)

var opts struct {
	MongoDbURI string `long:"mongodb-uri" env:"APP_MONGODB_URI" required:"true" description:"mongodb uri"`
	AppDb      string `long:"appdb" env:"MONGO_APP_DB" required:"true" description:"mongodb app db"`
}

func main() {
	if _, err := flags.Parse(&opts); err != nil {
		panic("failed to parse args")
	}

}
