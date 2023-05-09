package main

import (
	"github.com/evgenytr/metrics.git/internal/handlers"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc(`/update/`, handlers.ProcessRequest)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		log.Fatalln(err)
		panic(err)
	}
}
