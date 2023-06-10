package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/evgenytr/metrics.git/internal/config"
	"github.com/evgenytr/metrics.git/internal/handlers"
	"github.com/evgenytr/metrics.git/internal/logging"
	"github.com/evgenytr/metrics.git/internal/router"
	"github.com/evgenytr/metrics.git/internal/storage/memstorage"
	"github.com/go-chi/chi/v5"
)

func main() {

	host, storeInterval, fileStoragePath, restore := config.GetServerConfig()

	fmt.Println(*host, *storeInterval, *fileStoragePath, *restore)

	storage := memstorage.NewStorage()
	storageHandler := handlers.NewStorageHandler(storage)

	logger, err := logging.NewLogger(logging.NewDevelopment)
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		err = logger.Sync()
		if err != nil {
			fmt.Println(err)
		}
	}()
	sugar := logger.Sugar()

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

	r := router.Router(sugar, storageHandler)
	go listenAndServe(ctx, host, r)

	for {
		<-ctx.Done()
		err := context.Cause(ctx)
		if err != nil {
			log.Fatalln(err)
		}
		if *fileStoragePath != "" {
			err = storage.StoreMetrics(fileStoragePath)
			if err != nil {
				log.Print(err)
			}
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

func storeMetrics(ctx context.Context, storeInterval *time.Duration, fileStoragePath *string, storage memstorage.Storage) {
	_, cancelCtx := context.WithCancelCause(ctx)
	for {
		time.Sleep(*storeInterval)
		err := storage.StoreMetrics(fileStoragePath)

		if err != nil {
			fmt.Println("store metrics err")
			cancelCtx(err)
			return
		}
	}
}
