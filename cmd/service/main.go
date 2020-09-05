package main

import (
	"fmt"
	"github.com/misgorod/antibot/internal/config"
	"github.com/misgorod/antibot/internal/handler/limiter"
	"github.com/misgorod/antibot/internal/storage/zookeeper"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

func main() {
	config := config.Config{}
	err := envconfig.Process("antibot", &config)
	if err != nil {
		panic(err)
	}

	logger := logrus.New()
	if config.Trace {
		logger.SetLevel(logrus.TraceLevel)
	} else if config.Debug {
		logger.SetLevel(logrus.DebugLevel)
	}

	port := fmt.Sprintf(":%d", config.Port)

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(10 * time.Second))

	storage := zookeeper.New(config.ZkHosts)

	limiterHandler := limiter.New(storage, logger, config)

	router.Get("/api", limiterHandler.Handle)

	_ = http.ListenAndServe(port, router)
}
