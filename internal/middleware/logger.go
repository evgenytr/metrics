package middleware

import (
	"go.uber.org/zap"
	"net/http"
	"time"
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

			next.ServeHTTP(&lw, req)
			duration := time.Since(start)
			logger.Infoln(
				"uri", uri,
				"method", method,
				"duration", duration,
				"status", responseData.status,
				"size", responseData.size,
			)
		}
		return http.HandlerFunc(logFn)
	}
}
