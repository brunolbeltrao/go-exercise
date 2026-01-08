package main

import (
	"log"
	"net/http"
	"time"

	"go-exercise/internal/cache"
	"go-exercise/internal/config"
	httpapi "go-exercise/internal/http"
	"go-exercise/internal/http/handlers"
	"go-exercise/internal/kraken"
	"go-exercise/internal/ltp"
)

func main() {
	cfg := config.FromEnv()

	// Wire dependencies
	c := cache.NewMemoryCache(cfg.CacheTTL)
	k := kraken.NewRealClient(cfg.HTTPTimeout)
	svc := ltp.NewService(c, k)
	ltpHandler := handlers.NewLTPHandler(svc)
	router := httpapi.NewRouter(ltpHandler)

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("listening on :%s", cfg.Port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}
