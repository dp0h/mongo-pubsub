package rest

import (
	"github.com/dp0h/mongo-pubsub/pubsub"
)

type Srv struct {
	Router *chi.Mux
	pubSub *pubsub.PubSub
}

func NewRestSrv(pubSub *pubsub.PubSub) (*Srv, error) {
	s := Srv{
		pubSub: pubSub,
	}

	s.Router = chi.NewRouter()
	s.Router.Use(middleware.NoCache)

	s.Get("/ping", s.Ping)

	s.Post("/event", s.PostEvent)
	s.Get("/events", s.GetEvents)

	return &s, nil
}
