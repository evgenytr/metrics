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

func (r *loggingResponseWriter) Write(b []byte) (size int, err error) {
	size, err = r.ResponseWriter.Write(b)
	r.responseData.size += size
	return
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}
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
			)
		}
		return http.HandlerFunc(logFn)
	}
}
