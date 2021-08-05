package rest

import (
	"fmt"
	"github.com/dp0h/mongo-pubsub/pubsub"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/rs/zerolog/log"
	"net/http"
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

	s.Router.Get("/ping", s.Ping)

	s.Router.Post("/event", s.PostEvent)
	s.Router.Get("/events", s.GetEvents)

	return &s, nil
}

func (s *Srv) ListenAndServe(host string, port int) error {
	url := fmt.Sprintf("%s:%d", host, port)
	log.Info().Str("host", url).Msg("Started server")
	return http.ListenAndServe(url, s.Router)
}
