// Package router contains creation and configuration of new chi Router for metrics server service.
package router

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net"
	"os"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	"github.com/evgenytr/metrics.git/internal/handlers"
	"github.com/evgenytr/metrics.git/internal/middleware"
)

const CompressionLevel = 5

// Router method creates new Router.
func Router(sugar *zap.SugaredLogger, h *handlers.StorageHandler, key, cryptoKey, trustedSubnet string) (r *chi.Mux, err error) {
	r = chi.NewRouter()

	if trustedSubnet != "" {
		fmt.Println("middleware with trusted subnet used")
		_, ipNet, err := net.ParseCIDR(trustedSubnet)
		if err != nil {
			return nil, fmt.Errorf("failed to parse trusted subnet CIDR: %w", err)
		}
		withSubnetCheck := middleware.WithSubnetCheck(ipNet)
		r.Use(withSubnetCheck)
	}

	if key != "" {
		fmt.Println("middleware with signature check used")
		withSignatureCheck := middleware.WithSignatureCheck(key)
		r.Use(withSignatureCheck)
	}

	if cryptoKey != "" {
		dat, err := os.ReadFile(cryptoKey)
		if err != nil {
			return nil, fmt.Errorf("failed to read private key from file: %w", err)
		}

		block, _ := pem.Decode(dat)

		privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("failed to parse private key: %w", err)
		}

		fmt.Println("middleware with decryption is used")
		withDecryption := middleware.WithDecryption(privateKey)
		r.Use(withDecryption)
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

	return
}
