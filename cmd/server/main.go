package main

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"github.com/evgenytr/metrics.git/internal/config"
	"github.com/evgenytr/metrics.git/internal/handlers"
	"github.com/evgenytr/metrics.git/internal/middleware"
	"github.com/evgenytr/metrics.git/internal/storage/memstorage"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
	"log"
	"net/http"
)

func main() {

	host := getFlags()
	var cfg config.ServerConfig
	_ = env.Parse(&cfg)

	flag.Parse()

	if cfg.Host != "" {
		host = &cfg.Host
	}

	storage := memstorage.NewStorage()
	h := handlers.NewBaseHandler(storage)

	//TODO move logger to logging package
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalln(err)
	}
	defer logger.Sync()

	sugar := logger.Sugar()

	withLogging := middleware.WithLogging(sugar)

	r := chi.NewRouter()

	r.Use(withLogging)
	r.Use(chiMiddleware.AllowContentEncoding("gzip"))
	r.Use(chiMiddleware.Compress(5, "text/html", "application/json"))
	r.Post("/update/", h.ProcessPostUpdateJSONRequest)
	r.Post("/value/", h.ProcessPostValueJSONRequest)

	r.Post("/update/{type}/{name}/{value}", h.ProcessPostUpdateRequest)
	r.Get("/value/{type}/{name}", h.ProcessGetValueRequest)

	r.Get("/", h.ProcessGetListRequest)

	err = http.ListenAndServe(*host, r)
	if err != nil {
		log.Fatalln(err)
	}
}

func getFlags() (host *string) {
	host = flag.String("a", "localhost:8080", "host address")
	return
}
