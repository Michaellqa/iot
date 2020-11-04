package web

import (
	"context"
	"encoding/json"
	"github.com/Michaellqa/iot/app"
	"log"
	"net/http"
)

type AppServer struct {
	http.Server
}

func NewServer(ctx context.Context) http.Server {
	h := newHandler(ctx)

	mux := http.NewServeMux()
	mux.Handle("/", h)

	srv := http.Server{
		Addr:    ":8081",
		Handler: mux,
	}
	return srv
}

func newHandler(ctx context.Context) *handler {
	return &handler{
		globalCtx: ctx,
		busy:      make(chan struct{}, 1),
	}
}

type handler struct {
	globalCtx context.Context
	busy      chan struct{}
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	cfg := app.Config{}
	if err := json.NewDecoder(r.Body).Decode(&cfg); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	select {
	case h.busy <- struct{}{}:
	default:
		w.WriteHeader(http.StatusTooManyRequests)
		return
	}

	sup := app.Supervisor{}
	sup.Init(cfg)
	go func() {
		sup.Start(h.globalCtx)
		<-h.busy
	}()
}
