package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

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