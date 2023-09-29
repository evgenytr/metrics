// Package middleware contains custom middleware for metrics server service.
package middleware

import (
	"bytes"
	cryptoRand "crypto/rand"
	"crypto/rsa"
	"fmt"
	"hash"
	"io"
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
func WithDecryption(cryptoKey *rsa.PrivateKey) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		decryptFn := func(res http.ResponseWriter, req *http.Request) {
			encrypted := req.Header["Rsa-Encrypted"]

			if len(encrypted) != 0 && encrypted[0] == "true" {
				fmt.Println("encrypted header set in request")

				bodyBytes, err := io.ReadAll(req.Body)

				if err != nil {
					fmt.Println("err reading request body", err)
				}

				msgLen := len(bodyBytes)
				step := cryptoKey.PublicKey.Size()
				var decryptedBytes []byte

				for start := 0; start < msgLen; start += step {
					finish := start + step
					if finish > msgLen {
						finish = msgLen
					}

					decryptedBlockBytes, err := rsa.DecryptPKCS1v15(cryptoRand.Reader, cryptoKey, bodyBytes[start:finish])

					if err != nil {
						fmt.Println("err decrypting request body", err)
					}

					decryptedBytes = append(decryptedBytes, decryptedBlockBytes...)
				}

				req.Body = io.NopCloser(bytes.NewBuffer(decryptedBytes))
			}

			next.ServeHTTP(res, req)

		}
		return http.HandlerFunc(decryptFn)
	}
}
