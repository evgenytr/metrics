package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/evgenytr/metrics.git/internal/metric"
	"github.com/evgenytr/metrics.git/internal/storage/memstorage"
	"net/http"
	"net/url"
	"strings"
)

type BaseHandler struct {
	storage memstorage.Storage
}

func NewBaseHandler(storage memstorage.Storage) *BaseHandler {
	return &BaseHandler{
		storage: storage,
	}
}

func (h *BaseHandler) ProcessPostUpdateJSONRequest(res http.ResponseWriter, req *http.Request) {

	res.Header().Set("Content-Type", "application/json")

	dec := json.NewDecoder(req.Body)
	var currMetric metric.Metrics

	err := dec.Decode(&currMetric)
	if err != nil {
		processBadRequest(res, err)
		return
	}

	switch currMetric.MType {
	case "gauge":
		currMetric.Value, err = h.storage.UpdateGauge(currMetric.ID, currMetric.Value)
	case "counter":
		currMetric.Delta, err = h.storage.UpdateCounter(currMetric.ID, currMetric.Delta)
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

func (h *BaseHandler) ProcessPostValueJSONRequest(res http.ResponseWriter, req *http.Request) {

	res.Header().Set("Content-Type", "application/json")

	dec := json.NewDecoder(req.Body)
	var currMetric metric.Metrics

	err := dec.Decode(&currMetric)
	if err != nil {
		processBadRequest(res, err)
		return
	}

	fmt.Println(currMetric)

	switch currMetric.MType {
	case "gauge":
		currMetric.Value, err = h.storage.GetGaugeValue(currMetric.ID)
	case "counter":
		currMetric.Delta, err = h.storage.GetCounterValue(currMetric.ID)
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
func (h *BaseHandler) ProcessPostUpdateRequest(res http.ResponseWriter, req *http.Request) {

	const requiredRequestPathChunks = 5

	res.Header().Set("Content-Type", "text/html")

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

	_, err = h.storage.Update(metricType, metricName, metricValue)
	fmt.Println(err)
	if err != nil {
		processBadRequest(res, err)
		return
	}
	res.WriteHeader(http.StatusOK)
}

func (h *BaseHandler) ProcessGetValueRequest(res http.ResponseWriter, req *http.Request) {

	res.Header().Set("Content-Type", "text/html")

	const requiredRequestPathChunks = 4

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

	value, err := h.storage.ReadValue(metricType, metricName)

	if err != nil {
		res.WriteHeader(http.StatusNotFound)
		return
	}
	res.Write([]byte(value))
	res.WriteHeader(http.StatusOK)
}

func (h *BaseHandler) ProcessGetListRequest(res http.ResponseWriter, req *http.Request) {

	res.Header().Set("Content-Type", "text/html")

	metricsMap, err := h.storage.ListAll()

	if err != nil {
		processBadRequest(res, err)
		return
	}
	for key, value := range metricsMap {
		res.Write([]byte(fmt.Sprintf("%v\t%v\r", key, value)))
	}

	res.WriteHeader(http.StatusOK)
}

func processBadRequest(res http.ResponseWriter, err error) {
	res.Header().Set("Content-Type", "text/html")
	res.WriteHeader(http.StatusBadRequest)
	res.Write([]byte(fmt.Sprintf("Bad request, error %v", err)))
	fmt.Println(err)
}
