package rest

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/dp0h/mongo-pubsub/pubsub"
	"github.com/rs/zerolog/log"
)

func (s *Srv) Ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

type PostEventReq struct {
	Message string `json:"msg"`
}

func (s *Srv) PostEvent(w http.ResponseWriter, r *http.Request) {
	var postEventReq PostEventReq
	err := json.NewDecoder(r.Body).Decode(&postEventReq)
	if err != nil {
		log.Info().Err(err).Msg("Error parsing request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	topic := "v1_mongo_pubsub_event_post"
	event := &pubsub.Event{
		ID:      pubsub.Global.CounterMap[topic] + 1,
		Message: postEventReq.Message,
		Time:    time.Now().UTC(),
	}
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()

	err = s.pubSub.Publishers()[topic].Push(ctx, event)
	if err != nil {
		log.Info().Err(err).Msg("Error while adding event")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
}

func (s *Srv) GetEvents(w http.ResponseWriter, r *http.Request) {
	topic := "v1_mongo_pubsub_events_get"
	events, err := s.pubSub.GetEvents(topic)
	if err != nil {
		log.Warn().Err(err).Msg("Error while getting events")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(events); err != nil {
		log.Warn().Err(err).Msg("Error sending response")
	}
}
