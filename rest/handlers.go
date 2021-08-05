package rest

import (
	"encoding/json"
	"github.com/rs/zerolog/log"
	"net/http"
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

	err = s.pubSub.AddEvent(postEventReq.Message)
	if err != nil {
		log.Info().Err(err).Msg("Error while adding event")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
}

func (s *Srv) GetEvents(w http.ResponseWriter, r *http.Request) {
	events, err := s.pubSub.GetEvents()
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
