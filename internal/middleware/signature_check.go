package middleware

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"hash"
	"io"
	"net/http"
)

type (
	signatureResponseData struct {
		h hash.Hash
	}

	signatureResponseWriter struct {
		http.ResponseWriter
		responseData *signatureResponseData
	}
)

func (r *signatureResponseWriter) Write(b []byte) (size int, err error) {
	r.responseData.h.Write(b)
	size, err = r.ResponseWriter.Write(b)
	return
}

func WithSignatureCheck(key *string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		checkFn := func(res http.ResponseWriter, req *http.Request) {
			fmt.Println("WithSignatureCheck")
			responseData := &signatureResponseData{h: sha256.New()}
			sr := signatureResponseWriter{
				ResponseWriter: res,
				responseData:   responseData,
			}

			//compare incoming hash with calculated
			if key != nil && *key != "" {
				sr.Header().Set("Trailer", "HashSHA256")

				incomingHash := req.Header["Hashsha256"]
				if len(incomingHash) != 0 {
					fmt.Println(incomingHash, *key)

					hash := sha256.New()

					bodyBytes, _ := io.ReadAll(req.Body)

					req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

					keyBytes := []byte(*key)

					src := append(bodyBytes, keyBytes...)
					hash.Write(src)

					dst := hash.Sum(nil)

					decodedIncomingHash, _ := base64.StdEncoding.DecodeString(incomingHash[0])

					if !bytes.Equal(dst, decodedIncomingHash) {
						sr.WriteHeader(http.StatusBadRequest)
					}

				}
			}

			next.ServeHTTP(&sr, req)

			if key != nil && *key != "" {
				keyBytes := []byte(*key)
				sr.responseData.h.Write(keyBytes)
				dst := sr.responseData.h.Sum(nil)
				encodedDst := base64.StdEncoding.EncodeToString(dst)
				sr.Header().Set("HashSHA256", encodedDst)
				fmt.Println("Set response sha256 to ", encodedDst)
			}

		}
		return http.HandlerFunc(checkFn)
	}
}
