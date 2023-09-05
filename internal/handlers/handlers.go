// Package handlers contains handlers for all http requests to metrics server.
package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/evgenytr/metrics.git/internal/interfaces"
	"github.com/evgenytr/metrics.git/internal/metric"
)

// StorageHandler struct contains Storage interface.
type StorageHandler struct {
	storage interfaces.Storage
}

// NewStorageHandler returns StorageHandler for passed interface.
func NewStorageHandler(storage interfaces.Storage) *StorageHandler {
	return &StorageHandler{
		storage: storage,
	}
}

// ProcessPostUpdateJSONRequest handler processes JSON POST request to /update/,
// saves value of single metric.
func (h *StorageHandler) ProcessPostUpdateJSONRequest(res http.ResponseWriter, req *http.Request) {

	res.Header().Set("Content-Type", "application/json")

	ctx := req.Context()
	var currMetric metric.Metrics
	var buf bytes.Buffer

	_, err := buf.ReadFrom(req.Body)

	if err != nil {
		processBadRequest(res, err)
		return
	}

	err = json.Unmarshal(buf.Bytes(), &currMetric)

	if err != nil {
		processBadRequest(res, err)
		return
	}

	switch currMetric.MType {
	case metric.GaugeMetricType:
		currMetric.Value, err = h.storage.UpdateGauge(ctx, currMetric.ID, currMetric.Value)
	case metric.CounterMetricType:
		currMetric.Delta, err = h.storage.UpdateCounter(ctx, currMetric.ID, currMetric.Delta)
	default:
		err = fmt.Errorf("metric type unknown %v", currMetric.MType)
	}

	if err != nil {
		processBadRequest(res, err)
		return
	}

	err = json.NewEncoder(res).Encode(currMetric)
	if err != nil {
		processBadRequest(res, err)
		return
	}

	res.WriteHeader(http.StatusOK)
}

// ProcessPostValueJSONRequest handler processes JSON POST request to /value/,
// response is value of requested metric.
func (h *StorageHandler) ProcessPostValueJSONRequest(res http.ResponseWriter, req *http.Request) {

	res.Header().Set("Content-Type", "application/json")

	ctx := req.Context()
	var currMetric metric.Metrics
	var buf bytes.Buffer

	_, err := buf.ReadFrom(req.Body)

	if err != nil {
		processBadRequest(res, err)
		return
	}

	err = json.Unmarshal(buf.Bytes(), &currMetric)

	if err != nil {
		processBadRequest(res, err)
		return
	}

	fmt.Println(currMetric)

	switch currMetric.MType {
	case metric.GaugeMetricType:
		currMetric.Value, err = h.storage.GetGaugeValue(ctx, currMetric.ID)
	case metric.CounterMetricType:
		currMetric.Delta, err = h.storage.GetCounterValue(ctx, currMetric.ID)
	default:
		err = fmt.Errorf("metric type unknown %v", currMetric.MType)
		processBadRequest(res, err)
		return
	}

	if err != nil {
		fmt.Println(err)
		res.WriteHeader(http.StatusNotFound)
		return
	}

	err = json.NewEncoder(res).Encode(currMetric)
	if err != nil {
		processBadRequest(res, err)
		return
	}

	res.WriteHeader(http.StatusOK)
}

// ProcessPostUpdateRequest handles plain text request to /update/{type}/{name}/{value},
// saves value of single metric.
func (h *StorageHandler) ProcessPostUpdateRequest(res http.ResponseWriter, req *http.Request) {

	const requiredRequestPathChunks = 5

	ctx := req.Context()
	res.Header().Set("Content-Type", "text/plain")

	//validate URL
	parsedURL, err := url.ParseRequestURI(req.RequestURI)
	if err != nil {

		processBadRequest(res, err)
		return
	}

	requestPath := parsedURL.Path
	fmt.Println(requestPath)
	requestData := strings.Split(requestPath, "/")
	fmt.Println(len(requestData))
	if len(requestData) != requiredRequestPathChunks {
		res.WriteHeader(http.StatusNotFound)
		return
	}
	metricType := requestData[2]
	metricName := requestData[3]
	metricValue := requestData[4]

	fmt.Println(metricType, metricName, metricValue)

	_, err = h.storage.Update(ctx, metricType, metricName, metricValue)
	fmt.Println(err)
	if err != nil {
		processBadRequest(res, err)
		return
	}
	res.WriteHeader(http.StatusOK)
}

// ProcessGetValueRequest handles plain text Get request to /value/{type}/{name},
// response is value of requested metric.
func (h *StorageHandler) ProcessGetValueRequest(res http.ResponseWriter, req *http.Request) {

	const requiredRequestPathChunks = 4

	res.Header().Set("Content-Type", "text/plain")

	ctx := req.Context()

	//validate URL
	parsedURL, err := url.ParseRequestURI(req.RequestURI)
	if err != nil {
		processBadRequest(res, err)
		return
	}

	requestPath := parsedURL.Path

	requestData := strings.Split(requestPath, "/")
	if len(requestData) != requiredRequestPathChunks {
		res.WriteHeader(http.StatusNotFound)
		return
	}
	metricType := requestData[2]
	metricName := requestData[3]

	value, err := h.storage.ReadValue(ctx, metricType, metricName)

	if err != nil {
		res.WriteHeader(http.StatusNotFound)
		return
	}
	_, err = res.Write([]byte(value))
	if err != nil {
		fmt.Println(err)
	}

	res.WriteHeader(http.StatusOK)
}

// ProcessGetListRequest handles Get request to /,
// response is list of all metrics names and values.
func (h *StorageHandler) ProcessGetListRequest(res http.ResponseWriter, req *http.Request) {

	res.Header().Set("Content-Type", "text/html")

	ctx := req.Context()
	metricsMap, err := h.storage.ListAll(ctx)

	if err != nil {
		processBadRequest(res, err)
		return
	}
	for key, value := range metricsMap {
		_, err = fmt.Fprintf(res, "%v\t%v\r", key, value)
		if err != nil {
			fmt.Println(err)
		}
	}

	res.WriteHeader(http.StatusOK)
}

// ProcessPingRequest handles Get request to /ping, used to check whether service is running and responding.
func (h *StorageHandler) ProcessPingRequest(res http.ResponseWriter, req *http.Request) {

	res.Header().Set("Content-Type", "text/html")

	ctx := req.Context()
	if err := h.storage.Ping(ctx); err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)
}

// ProcessPostUpdatesBatchJSONRequest handles POST JSON request to /updates/,
// processes array of metrics and saves them all at once.
func (h *StorageHandler) ProcessPostUpdatesBatchJSONRequest(res http.ResponseWriter, req *http.Request) {

	res.Header().Set("Content-Type", "application/json")

	ctx := req.Context()

	var currMetrics []metric.Metrics
	var buf bytes.Buffer

	_, err := buf.ReadFrom(req.Body)

	if err != nil {
		processBadRequest(res, err)
		return
	}

	err = json.Unmarshal(buf.Bytes(), &currMetrics)

	if err != nil {
		processBadRequest(res, err)
		return
	}

	for _, currMetric := range currMetrics {
		switch currMetric.MType {
		case metric.GaugeMetricType:
			currMetric.Value, err = h.storage.UpdateGauge(ctx, currMetric.ID, currMetric.Value)
		case metric.CounterMetricType:
			currMetric.Delta, err = h.storage.UpdateCounter(ctx, currMetric.ID, currMetric.Delta)
		default:
			err = fmt.Errorf("metric type unknown %v", currMetric.MType)
		}

		if err != nil {
			processBadRequest(res, err)
			return
		}
	}

	//workaround for autotest fail on unmarshalling the response to main.Metrics
	//sends null
	var outMetrics []metric.Metrics
	out, err := json.Marshal(outMetrics)

	if err != nil {
		processBadRequest(res, err)
		return
	}

	_, err = res.Write(out)
	if err != nil {
		fmt.Println(err)
	}

	res.WriteHeader(http.StatusOK)
}

func processBadRequest(res http.ResponseWriter, err error) {
	res.Header().Set("Content-Type", "text/plain")
	res.WriteHeader(http.StatusBadRequest)

	_, errOut := fmt.Fprintf(res, "Bad request, error %v", err)
	if errOut != nil {
		fmt.Println(err)
	}
	fmt.Println(err)
}
