package router

import (
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	"github.com/evgenytr/metrics.git/internal/handlers"
	"github.com/evgenytr/metrics.git/internal/middleware"
)

const CompressionLevel = 5

func Router(sugar *zap.SugaredLogger, h *handlers.StorageHandler) *chi.Mux {
	r := chi.NewRouter()

	withLogging := middleware.WithLogging(sugar)

	r.Use(withLogging)
	r.Use(chiMiddleware.AllowContentEncoding("gzip"))
	r.Use(chiMiddleware.Compress(CompressionLevel, "text/html", "application/json"))

	r.Get("/ping", h.ProcessPingRequest)

	r.Post("/updates/", h.ProcessPostUpdatesBatchJSONRequest)

	r.Post("/update/", h.ProcessPostUpdateJSONRequest)
	r.Post("/value/", h.ProcessPostValueJSONRequest)

	r.Post("/update/{type}/{name}/{value}", h.ProcessPostUpdateRequest)
	r.Get("/value/{type}/{name}", h.ProcessGetValueRequest)

	r.Get("/", h.ProcessGetListRequest)

	return r
}
