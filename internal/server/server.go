package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/EwvwGeN/EffectiveMobile_assignment/internal/config"
	"github.com/gorilla/mux"
)

type server struct {
	cfg config.HttpConfig
	log *slog.Logger
	router *mux.Router
}

func NewHttpServer(cfg config.HttpConfig, log *slog.Logger) *server {
	return &server{
		cfg: cfg,
		log: log,
		router: mux.NewRouter(),
	}
}

func (s *server) RunServer(ctx context.Context) (errCloseCh chan error) {
	s.log.Info("starting server")
	errCloseCh = make(chan error)
	srv := &http.Server{
		Handler: s.router,
		Addr:    fmt.Sprintf("%s:%s", s.cfg.Host, s.cfg.Port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	s.log.Info("starting listening", slog.String("addres", srv.Addr))
	go func() {
		<-ctx.Done()
		s.log.Info("Graceful shutdown http server")
		errCloseCh <- srv.Shutdown(context.Background())
	}()
	go srv.ListenAndServe()
	return
}

func(s *server) RegisterHandler(url string, handler http.HandlerFunc, method string) {
	s.router.HandleFunc(
		url,
		handler,
	).Methods(method)
}