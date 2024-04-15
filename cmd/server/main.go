package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	v1 "github.com/EwvwGeN/EffectiveMobile_assignment/http/v1"
	c "github.com/EwvwGeN/EffectiveMobile_assignment/internal/config"
	l "github.com/EwvwGeN/EffectiveMobile_assignment/internal/logger"
	"github.com/EwvwGeN/EffectiveMobile_assignment/internal/server"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config", "", "path to config file")
}
func main() {
	flag.Parse()
	cfg, err := c.LoadConfig(configPath)
	if err != nil {
		panic(fmt.Sprintf("cant load config from path %s: %s", configPath, err.Error()))
	}
	logger := l.SetupLogger(cfg.LogLevel)
	logger.Info("logger is initiated")
	logger.Debug("config data", slog.Any("config", cfg))
	mainCtx, cancel := context.WithCancel(context.Background())

	hserver := server.NewHttpServer(cfg.HttpConfig, logger)

	hserver.RegisterHandler(
		"/api/cars/add",
		v1.CarAdd(logger, nil),
		http.MethodPost,
	)
	hserver.RegisterHandler(
		"/api/car/{carId}",
		v1.CarGetOne(logger, nil),
		http.MethodGet,
	)
	hserver.RegisterHandler(
		"/api/cars",
		v1.CarGetAll(logger, nil, nil),
		http.MethodGet,
	)
	hserver.RegisterHandler(
		"/api/car/{carId}/edit",
		v1.CarEdit(logger, cfg.ValidatorConfig, nil),
		http.MethodPost,
	)
	hserver.RegisterHandler(
		"/api/car/{carId}/delete",
		v1.CarDelete(logger, nil),
		http.MethodPost,
	)

	logger.Info("loading end")

	errCh := hserver.RunServer(mainCtx)
	stopChecker := make(chan os.Signal, 1)
	signal.Notify(stopChecker, syscall.SIGTERM, syscall.SIGINT)
	<- stopChecker
	logger.Info("stopping service")
	cancel()
	err = <-errCh
	if err != nil {
		logger.Error("error while stopping http server", slog.String("error", err.Error()))
	}
	logger.Info("service stoped successfully")
}