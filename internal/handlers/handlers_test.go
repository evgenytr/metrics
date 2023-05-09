package handlers

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestProcessRequest(t *testing.T) {
	type want struct {
		code        int
		contentType string
	}
	tests := []struct {
		name    string
		method  string
		request string
		want    want
	}{
		{
			name:    "Positive new gauge metric",
			method:  http.MethodPost,
			request: "/update/gauge/gaugeMetric/123",
			want: want{
				code:        http.StatusOK,
				contentType: "text/plain",
			},
		},
		{
			name:    "Positive update gauge metric",
			method:  http.MethodPost,
			request: "/update/gauge/gaugeMetric/321",
			want: want{
				code:        http.StatusOK,
				contentType: "text/plain",
			},
		},
		{
			name:    "Negative update gauge metric as counter",
			method:  http.MethodPost,
			request: "/update/counter/gaugeMetric/321",
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/plain",
			},
		},
		{
			name:    "Positive new counter metric",
			method:  http.MethodPost,
			request: "/update/counter/counterMetric/123",
			want: want{
				code:        http.StatusOK,
				contentType: "text/plain",
			},
		},
		{
			name:    "Positive update counter metric",
			method:  http.MethodPost,
			request: "/update/counter/counterMetric/123",
			want: want{
				code:        http.StatusOK,
				contentType: "text/plain",
			},
		},
		{
			name:    "Negative update counter metric as gauge",
			method:  http.MethodPost,
			request: "/update/gauge/counterMetric/321",
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/plain",
			},
		},
		{
			name:    "Negative update metric with negative value",
			method:  http.MethodPost,
			request: "/update/gauge/gaugeMetric/-321",
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/plain",
			},
		},
		{
			name:    "Negative update metric without value",
			method:  http.MethodPost,
			request: "/update/gauge/gaugeMetric",
			want: want{
				code:        http.StatusNotFound,
				contentType: "text/plain",
			},
		},
		{
			name:    "Negative update metric without name",
			method:  http.MethodPost,
			request: "/update/gauge",
			want: want{
				code:        http.StatusNotFound,
				contentType: "text/plain",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, tt.request, nil)
			w := httptest.NewRecorder()
			ProcessRequest(w, request)
			res := w.Result()
			assert.Equal(t, tt.want.code, res.StatusCode)
			assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))
			res.Body.Close()
		})
	}
}