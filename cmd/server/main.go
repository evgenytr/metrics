// Package main initializes metrics server service and starts it.
package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/evgenytr/metrics.git/internal/config"
	errorHandling "github.com/evgenytr/metrics.git/internal/errors"
	"github.com/evgenytr/metrics.git/internal/handlers"
	"github.com/evgenytr/metrics.git/internal/interfaces"
	"github.com/evgenytr/metrics.git/internal/logging"
	"github.com/evgenytr/metrics.git/internal/router"
	"github.com/evgenytr/metrics.git/internal/storage"
	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {

	var err error
	var db *sql.DB

	host, storeInterval, fileStoragePath, restore, dbDSN, key := config.GetServerConfig()

	fmt.Println(*host, *storeInterval, *fileStoragePath, *restore, *dbDSN, *key)

	if *dbDSN != "" {
		db, err = sql.Open("pgx", *dbDSN)

		if err != nil {
			log.Fatalln(err)
		}

		defer func() {
			err = db.Close()
			if err != nil {
				fmt.Println(err)
			}
		}()

		//ping
		err = db.Ping()
		if err != nil {
			log.Fatalln(err)
		}
	}

	useDatabase := db != nil
	appStorage := storage.NewStorage(db, fileStoragePath)

	storageHandler := handlers.NewStorageHandler(appStorage)

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

	err = appStorage.InitializeMetrics(ctx, restore)
	if err != nil {
		if useDatabase {
			log.Fatalln(err)
		}
		//error is not critical unless it's database storage
		fmt.Println(err)
	}

	if *storeInterval != 0 {
		go storeMetrics(ctx, storeInterval, appStorage)
	}

	r := router.Router(sugar, storageHandler, key)
	go listenAndServe(ctx, host, r)

	for {
		<-ctx.Done()
		err := context.Cause(ctx)
		if err != nil {
			log.Fatalln(err)
		}

		err = appStorage.StoreMetrics(ctx)
		if err != nil {
			log.Print(err)
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

func storeMetrics(ctx context.Context, storeInterval *time.Duration, storage interfaces.Storage) {
	_, cancelCtx := context.WithCancelCause(ctx)
	for {
		time.Sleep(*storeInterval)
		err := storage.StoreMetrics(ctx)

		if err != nil {
			for _, retryInterval := range errorHandling.RepeatedAttemptsIntervals {
				time.Sleep(*retryInterval)
				err = storage.StoreMetrics(ctx)
				if err == nil {
					break
				}
			}
		}

		if err != nil {
			fmt.Println("store metrics err ", err)
			cancelCtx(err)
			return
		}
	}
}
