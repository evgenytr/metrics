package handlers

import (
	"fmt"
	"github.com/evgenytr/metrics.git/internal/storage/memstorage"
	"net/http"
	"strings"
)

type Storage interface {
	Update(metricType, name, value string) error
	ReadValue(metricType, name string) (string, error)
	ListAll() (map[string]string, error)
}

var (
	storage Storage = memstorage.GetNewStorage()
)

func ProcessPostUpdateRequest(res http.ResponseWriter, req *http.Request) {

	res.Header().Set("Content-Type", "text/plain")
	requestData := strings.Split(req.RequestURI[1:], "/")
	if len(requestData) != 4 {
		res.WriteHeader(http.StatusNotFound)
		return
	}
	metricType := requestData[1]
	metricName := requestData[2]
	metricValue := requestData[3]

	fmt.Println(metricType, metricName, metricValue)

	err := storage.Update(metricType, metricName, metricValue)
	fmt.Println(err)
	if err != nil {
		processBadRequest(res)
		return
	}
	res.WriteHeader(http.StatusOK)
}

func ProcessGetValueRequest(res http.ResponseWriter, req *http.Request) {

	res.Header().Set("Content-Type", "text/plain")

	requestData := strings.Split(req.RequestURI[1:], "/")
	if len(requestData) != 3 {
		res.WriteHeader(http.StatusNotFound)
		return
	}
	metricType := requestData[1]
	metricName := requestData[2]

	value, err := storage.ReadValue(metricType, metricName)
	fmt.Println(err)
	if err != nil {
		res.WriteHeader(http.StatusNotFound)
		return
	}
	res.Write([]byte(value))
	res.WriteHeader(http.StatusOK)
}

func ProcessGetListRequest(res http.ResponseWriter, req *http.Request) {

	res.Header().Set("Content-Type", "text/plain")

	metricsMap, err := storage.ListAll()
	fmt.Println(err)
	if err != nil {
		processBadRequest(res)
		return
	}
	for key, value := range metricsMap {
		res.Write([]byte(fmt.Sprintf("%v\t%v\r", key, value)))
	}

	res.WriteHeader(http.StatusOK)
}

func processBadRequest(res http.ResponseWriter) {
	res.Header().Set("Content-Type", "text/plain")
	res.WriteHeader(http.StatusBadRequest)
	res.Write([]byte("Bad request"))
}
