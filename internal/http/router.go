package httpapi

import (
	"net/http"

	"go-exercise/internal/http/handlers"
	"go-exercise/internal/http/middleware"
)

func NewRouter(ltpHandler http.Handler) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/api/v1/ltp", ltpHandler)
	mux.Handle("/healthz", handlers.NewHealthHandler())
	mux.Handle("/readyz", handlers.NewReadyHandler())
	return middleware.Logging(mux)
}
