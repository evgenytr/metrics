package handlers

import (
	"fmt"
	"github.com/evgenytr/metrics.git/internal/storage/memstorage"
	"net/http"
	"strings"
)

type Storage interface {
	Update(metricType, name, value string) error
}

var (
	storage Storage = memstorage.GetNewStorage()
)

func ProcessRequest(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		processUpdateMetricPostRequest(res, req)
		return
	}

	processBadRequest(res)
}

func processUpdateMetricPostRequest(res http.ResponseWriter, req *http.Request) {

	requestData := strings.Split(req.RequestURI[1:], "/")

	if len(requestData) != 4 {
		res.WriteHeader(http.StatusNotFound)
		return
	}
	metricType := requestData[1]
	metricName := requestData[2]
	metricValue := requestData[3]

	err := storage.Update(metricType, metricName, metricValue)
	fmt.Println(err)
	if err != nil {
		processBadRequest(res)
		return
	}
	fmt.Println(req.RequestURI)
	res.WriteHeader(http.StatusOK)
}

func processBadRequest(res http.ResponseWriter) {
	res.WriteHeader(http.StatusBadRequest)
	res.Write([]byte("Bad request"))
}
