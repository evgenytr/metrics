package main

import (
	"context"
	"fmt"
	"github.com/evgenytr/metrics.git/internal/config"
	"github.com/evgenytr/metrics.git/internal/handlers"
	"github.com/evgenytr/metrics.git/internal/middleware"
	"github.com/evgenytr/metrics.git/internal/storage/memstorage"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
	"log"
	"net/http"
	"time"
)

func main() {

	host, storeInterval, fileStoragePath, restore := config.GetServerConfig()

	fmt.Println(*host, *storeInterval, *fileStoragePath, *restore)

	storage := memstorage.NewStorage()
	h := handlers.NewBaseHandler(storage)

	//TODO move logger to logging package (?)
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalln(err)
	}
	defer logger.Sync()
	sugar := logger.Sugar()

	withLogging := middleware.WithLogging(sugar)

	ctx := context.Background()
	if *fileStoragePath != "" {
		if *restore {
			err = storage.LoadMetrics(fileStoragePath)
			//non fatal error
			if err != nil {
				log.Print(err)
			}
		}

		if *storeInterval != 0 {
			go storeMetrics(ctx, storeInterval, fileStoragePath, storage)
		}

	}
	r := chi.NewRouter()

	r.Use(withLogging)
	r.Use(chiMiddleware.AllowContentEncoding("gzip"))
	r.Use(chiMiddleware.Compress(5, "text/html", "application/json"))
	r.Post("/update/", h.ProcessPostUpdateJSONRequest)
	r.Post("/value/", h.ProcessPostValueJSONRequest)

	r.Post("/update/{type}/{name}/{value}", h.ProcessPostUpdateRequest)
	r.Get("/value/{type}/{name}", h.ProcessGetValueRequest)

	r.Get("/", h.ProcessGetListRequest)

	go listenAndServe(ctx, host, r)

	for {
		<-ctx.Done()
		err := context.Cause(ctx)
		if err != nil {
			log.Fatalln(err)
		}
		if *fileStoragePath != "" {
			storage.StoreMetrics(fileStoragePath)
		}
	}
}

func listenAndServe(ctx context.Context, host *string, r *chi.Mux) {
	_, cancelCtx := context.WithCancelCause(ctx)
	err := http.ListenAndServe(*host, r)
	if err != nil {
		fmt.Println("listenAndServe err", err)
		cancelCtx(err)
	}
}

func storeMetrics(ctx context.Context, storeInterval *float64, fileStoragePath *string, storage memstorage.Storage) {
	_, cancelCtx := context.WithCancelCause(ctx)
	for {
		time.Sleep(time.Duration(*storeInterval) * time.Second)
		err := storage.StoreMetrics(fileStoragePath)

		if err != nil {
			fmt.Println("store metrics err")
			cancelCtx(err)
			return
		}
	}
}
