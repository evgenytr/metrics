// Package router contains creation and configuration of new chi Router for metrics server service.
package router

import (
	"fmt"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	"github.com/evgenytr/metrics.git/internal/handlers"
	"github.com/evgenytr/metrics.git/internal/middleware"
)

const CompressionLevel = 5

// Router method creates new Router.
func Router(sugar *zap.SugaredLogger, h *handlers.StorageHandler, key *string) *chi.Mux {
	r := chi.NewRouter()

	if key != nil && *key != "" {
		fmt.Println("middleware with signature check used")
		withSignatureCheck := middleware.WithSignatureCheck(key)
		r.Use(withSignatureCheck)
	}

	withLogging := middleware.WithLogging(sugar)
	r.Use(withLogging)

	r.Use(chiMiddleware.AllowContentEncoding("gzip"))
	r.Use(chiMiddleware.Compress(CompressionLevel, "text/html", "application/json"))
	r.Mount("/debug", chiMiddleware.Profiler())

	r.Get("/ping", h.ProcessPingRequest)

	r.Post("/updates/", h.ProcessPostUpdatesBatchJSONRequest)

	r.Post("/update/", h.ProcessPostUpdateJSONRequest)
	r.Post("/value/", h.ProcessPostValueJSONRequest)

	r.Post("/update/{type}/{name}/{value}", h.ProcessPostUpdateRequest)
	r.Get("/value/{type}/{name}", h.ProcessGetValueRequest)

	r.Get("/", h.ProcessGetListRequest)

	return r
}
