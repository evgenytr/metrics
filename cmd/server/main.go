package main

import (
	"github.com/evgenytr/metrics.git/internal/handlers"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func main() {
	r := chi.NewRouter()
	r.Post("/update/{type}/{name}/{value}", handlers.ProcessPostUpdateRequest)
	r.Get("/value/{type}/{name}", handlers.ProcessGetValueRequest)
	r.Get("/", handlers.ProcessGetListRequest)

	err := http.ListenAndServe(`:8080`, r)
	if err != nil {
		log.Fatalln(err)
		panic(err)
	}
}
