package rest

import "net/http"

func (s *Srv) Ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (s *Srv) PostEvent(w http.ResponseWriter, r *http.Request) {
	//
}

func (s *Srv) GetEvents(w http.ResponseWriter, r *http.Request) {
	//
}
