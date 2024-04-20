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

	"github.com/EwvwGeN/EffectiveMobile_assignment/http/helper"
	"github.com/EwvwGeN/EffectiveMobile_assignment/http/parser"
	v1 "github.com/EwvwGeN/EffectiveMobile_assignment/http/v1"
	c "github.com/EwvwGeN/EffectiveMobile_assignment/internal/config"
	l "github.com/EwvwGeN/EffectiveMobile_assignment/internal/logger"
	"github.com/EwvwGeN/EffectiveMobile_assignment/internal/server"
	"github.com/EwvwGeN/EffectiveMobile_assignment/internal/service"
	"github.com/EwvwGeN/EffectiveMobile_assignment/internal/storage/postgres"

	_ "github.com/EwvwGeN/EffectiveMobile_assignment/http/swagger"

	httpSwagger "github.com/swaggo/http-swagger"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config", "", "path to config file")
}

// @title Swagger для микросервиса Cars
// @version 1.0
// @description Swagger для микросервиса Cars
//
// @host localhost:9099
// @BasePath /
//go:generate swag init --parseDependency --parseInternal  -d .,../../http/v1 -o ../../http/swagger/
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

	carInfoGetter, err := helper.GetCarInfoGetter(logger, cfg.CarInfoGetterUrl, parser.ParseFromExternalApi)
	if err != nil {
		logger.Error("failed to initialise info getter", slog.String("error", err.Error()))
		os.Exit(1)
	}

	postgresRepo, err := postgres.NewPostgresProvider(context.Background(), cfg.PostgresConfig)
	if err != nil {
		logger.Error("failed to initialise postgres provider", slog.String("error", err.Error()))
		os.Exit(1)
	}

	carService := service.NewCarService(logger, postgresRepo, carInfoGetter)

	hserver := server.NewHttpServer(cfg.HttpConfig, logger)

	hserver.RegisterHandler(
		"/api/cars/add",
		v1.CarAdd(logger, carService),
		http.MethodPost,
	)
	hserver.RegisterHandler(
		"/api/car/{carId}",
		v1.CarGetOne(logger, carService),
		http.MethodGet,
	)
	hserver.RegisterHandler(
		"/api/cars",
		v1.CarGetAll(logger, carService, postgres.AddFilter),
		http.MethodGet,
	)
	hserver.RegisterHandler(
		"/api/car/{carId}/edit",
		v1.CarEdit(logger, cfg.ValidatorConfig, carService),
		http.MethodPatch,
	)
	hserver.RegisterHandler(
		"/api/car/{carId}/delete",
		v1.CarDelete(logger, carService),
		http.MethodDelete,
	)
	swagParams := []func(*httpSwagger.Config){
		httpSwagger.URL("doc.json"),
	}
	hserver.RegisterHandler(
		"/api/swagger/{*}",
		httpSwagger.Handler(swagParams...),
		http.MethodGet,
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