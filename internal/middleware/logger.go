// Package middleware contains custom middleware for metrics server service.
package middleware

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

type (
	responseData struct {
		status int
		size   int
	}

	loggingResponseWriter struct {
		http.ResponseWriter
		responseData *responseData
	}
)

// Write is helper method to write bytes to response.
func (r *loggingResponseWriter) Write(b []byte) (size int, err error) {
	size, err = r.ResponseWriter.Write(b)
	r.responseData.size += size
	return
}

// WriteHeader is helper method to write response header.
func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

// WithLogging wraps handler with logging middleware.
func WithLogging(logger *zap.SugaredLogger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		logFn := func(res http.ResponseWriter, req *http.Request) {
			start := time.Now()
			responseData := &responseData{}
			lw := loggingResponseWriter{
				ResponseWriter: res,
				responseData:   responseData,
			}
			uri := req.RequestURI
			method := req.Method
			encoding := req.Header["Accept-Encoding"]
			hash := req.Header["Hashsha256"]
			encrypted := req.Header["Rsa-Encrypted"]
			realIP := req.Header["X-Real-Ip"]

			next.ServeHTTP(&lw, req)
			duration := time.Since(start)
			logger.Infoln(
				"uri", uri,
				"method", method,
				"encoding", encoding,
				"duration", duration,
				"status", responseData.status,
				"size", responseData.size,
				"hash", hash,
				"encrypted", encrypted,
				"request IP", realIP,
			)
		}
		return http.HandlerFunc(logFn)
	}
}
