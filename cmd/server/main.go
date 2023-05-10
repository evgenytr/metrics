package main

import (
	"flag"
	"github.com/evgenytr/metrics.git/internal/handlers"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

var (
	host *string
)

func init() {
	host = flag.String("a", "localhost:8080", "host address")
}

func main() {
	flag.Parse()

	r := chi.NewRouter()
	r.Post("/update/{type}/{name}/{value}", handlers.ProcessPostUpdateRequest)
	r.Get("/value/{type}/{name}", handlers.ProcessGetValueRequest)
	r.Get("/", handlers.ProcessGetListRequest)

	err := http.ListenAndServe(*host, r)
	if err != nil {
		log.Fatalln(err)
		panic(err)
	}
}
