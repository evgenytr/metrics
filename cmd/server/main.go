package main

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"github.com/evgenytr/metrics.git/internal/handlers"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

type Config struct {
	Host string `env:"ADDRESS"`
}

func main() {

	host := getFlags()
	var cfg Config
	_ = env.Parse(&cfg)

	flag.Parse()

	if cfg.Host != "" {
		host = &(cfg.Host)
	}

	r := chi.NewRouter()
	r.Post("/update/{type}/{name}/{value}", handlers.ProcessPostUpdateRequest)
	r.Get("/value/{type}/{name}", handlers.ProcessGetValueRequest)
	r.Get("/", handlers.ProcessGetListRequest)

	err := http.ListenAndServe(*host, r)
	if err != nil {
		log.Fatalln(err)
	}
}

func getFlags() (host *string) {
	host = flag.String("a", "localhost:8080", "host address")
	return
}
