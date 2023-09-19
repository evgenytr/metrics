// Package middleware contains custom middleware for metrics server service.
package middleware

import (
	"crypto/sha256"
	"fmt"
	"hash"
	"net/http"
)

type (
	decryptionResponseData struct {
		h hash.Hash
	}

	decryptionResponseWriter struct {
		http.ResponseWriter
		responseData *decryptionResponseData
	}
)

// Write is helper method to write bytes to response.
func (r *decryptionResponseWriter) Write(b []byte) (size int, err error) {
	r.responseData.h.Write(b)
	size, err = r.ResponseWriter.Write(b)
	return
}

// WithDecryption wraps handler with middleware to decrypt body.
func WithDecryption(cryptoKey string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		decryptFn := func(res http.ResponseWriter, req *http.Request) {

			responseData := &decryptionResponseData{h: sha256.New()}
			dr := decryptionResponseWriter{
				ResponseWriter: res,
				responseData:   responseData,
			}

			next.ServeHTTP(&dr, req)

			if cryptoKey != "" {
				fmt.Println(req.RequestURI)
				fmt.Println("use decrypt")
			}

		}
		return http.HandlerFunc(decryptFn)
	}
}
