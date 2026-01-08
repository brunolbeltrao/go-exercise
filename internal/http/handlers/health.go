package handlers

import "net/http"

type HealthHandler struct{}
type ReadyHandler struct{}

func NewHealthHandler() http.Handler {
	return HealthHandler{}
}

func NewReadyHandler() http.Handler {
	return ReadyHandler{}
}

func (HealthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (ReadyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
