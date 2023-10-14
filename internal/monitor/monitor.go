// Package monitor contains bulk of working code for metrics agent service.
package monitor

import (
	"bytes"
	"context"
	cryptoRand "crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/evgenytr/metrics.git/pkg/api/v1"

	"github.com/evgenytr/metrics.git/internal/metric"
	"github.com/go-resty/resty/v2"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	_ "github.com/shirou/gopsutil/v3/net"
)

type monitor struct {
	metrics     map[string]*metric.Metrics
	hostAddress string
	hostGrpc    string
	key         string
	cryptoKey   string
	wg          *sync.WaitGroup
}

// Monitor interface describes .
type Monitor interface {
	PollMetrics() error
	PollAdditionalMetrics() error
	ReportMetrics() error
	ReportMetricsGrpc() error
	ResetPollCount()
}

func initMap() (initialMap map[string]*metric.Metrics, err error) {

	var metricList = map[string]string{
		"Alloc":           metric.GaugeMetricType,
		"BuckHashSys":     metric.GaugeMetricType,
		"Frees":           metric.GaugeMetricType,
		"GCSys":           metric.GaugeMetricType,
		"HeapAlloc":       metric.GaugeMetricType,
		"HeapIdle":        metric.GaugeMetricType,
		"HeapInuse":       metric.GaugeMetricType,
		"HeapObjects":     metric.GaugeMetricType,
		"HeapReleased":    metric.GaugeMetricType,
		"HeapSys":         metric.GaugeMetricType,
		"LastGC":          metric.GaugeMetricType,
		"Lookups":         metric.GaugeMetricType,
		"MCacheInuse":     metric.GaugeMetricType,
		"MCacheSys":       metric.GaugeMetricType,
		"MSpanInuse":      metric.GaugeMetricType,
		"MSpanSys":        metric.GaugeMetricType,
		"Mallocs":         metric.GaugeMetricType,
		"NextGC":          metric.GaugeMetricType,
		"OtherSys":        metric.GaugeMetricType,
		"PauseTotalNs":    metric.GaugeMetricType,
		"StackInuse":      metric.GaugeMetricType,
		"StackSys":        metric.GaugeMetricType,
		"Sys":             metric.GaugeMetricType,
		"TotalAlloc":      metric.GaugeMetricType,
		"PollCount":       metric.CounterMetricType,
		"NumForcedGC":     metric.GaugeMetricType,
		"NumGC":           metric.GaugeMetricType,
		"GCCPUFraction":   metric.GaugeMetricType,
		"RandomValue":     metric.GaugeMetricType,
		"TotalMemory":     metric.GaugeMetricType,
		"FreeMemory":      metric.GaugeMetricType,
		"CPUutilization1": metric.GaugeMetricType,
	}

	initialMap = make(map[string]*metric.Metrics, len(metricList))
	for metricName, metricType := range metricList {
		switch metricType {
		case metric.GaugeMetricType:
			var value float64
			initialMap[metricName], err = metric.CreateGauge(metricName, value)
		case metric.CounterMetricType:
			var value int64
			initialMap[metricName], err = metric.CreateCounter(metricName, value)
		}
		if err != nil {
			return
		}
	}

	return
}

// NewMonitor returns
func NewMonitor(hostAddress, hostGrpc, key, cryptoKey string, wg *sync.WaitGroup) (m Monitor, err error) {
	initialMap, err := initMap()
	if err != nil {
		return
	}
	m = &monitor{
		metrics:     initialMap,
		hostAddress: hostAddress,
		hostGrpc:    hostGrpc,
		key:         key,
		cryptoKey:   cryptoKey,
		wg:          wg,
	}
	return
}

func (m *monitor) ResetPollCount() {
	_ = m.metrics["PollCount"].ResetCounter()
}

func updateGaugeMetric(metric *metric.Metrics, value float64) (err error) {
	_, err = metric.UpdateGauge(value)
	return
}

func updateCounterMetric(metric *metric.Metrics, value int64) (err error) {
	_, err = metric.UpdateCounter(value)
	return
}

func (m *monitor) PollMetrics() (err error) {
	fmt.Println("pollMetrics")
	var rtm runtime.MemStats
	runtime.ReadMemStats(&rtm)

	err = updateGaugeMetric(m.metrics["Alloc"], float64(rtm.Alloc))
	if err != nil {
		return
	}

	err = updateGaugeMetric(m.metrics["BuckHashSys"], float64(rtm.BuckHashSys))
	if err != nil {
		return
	}

	err = updateGaugeMetric(m.metrics["Frees"], float64(rtm.Frees))
	if err != nil {
		return
	}

	err = updateGaugeMetric(m.metrics["GCCPUFraction"], rtm.GCCPUFraction)
	if err != nil {
		return
	}

	err = updateGaugeMetric(m.metrics["GCSys"], float64(rtm.GCSys))
	if err != nil {
		return
	}

	err = updateGaugeMetric(m.metrics["HeapAlloc"], float64(rtm.HeapAlloc))
	if err != nil {
		return
	}

	err = updateGaugeMetric(m.metrics["HeapIdle"], float64(rtm.HeapIdle))
	if err != nil {
		return
	}

	err = updateGaugeMetric(m.metrics["HeapInuse"], float64(rtm.HeapInuse))
	if err != nil {
		return
	}

	err = updateGaugeMetric(m.metrics["HeapObjects"], float64(rtm.HeapObjects))
	if err != nil {
		return
	}

	err = updateGaugeMetric(m.metrics["HeapReleased"], float64(rtm.HeapReleased))
	if err != nil {
		return
	}

	err = updateGaugeMetric(m.metrics["HeapSys"], float64(rtm.HeapSys))
	if err != nil {
		return
	}

	err = updateGaugeMetric(m.metrics["LastGC"], float64(rtm.LastGC))
	if err != nil {
		return
	}

	err = updateGaugeMetric(m.metrics["Lookups"], float64(rtm.Lookups))
	if err != nil {
		return
	}

	err = updateGaugeMetric(m.metrics["MCacheInuse"], float64(rtm.MCacheInuse))
	if err != nil {
		return
	}

	err = updateGaugeMetric(m.metrics["MCacheSys"], float64(rtm.MCacheSys))
	if err != nil {
		return
	}

	err = updateGaugeMetric(m.metrics["MSpanInuse"], float64(rtm.MSpanInuse))
	if err != nil {
		return
	}

	err = updateGaugeMetric(m.metrics["MSpanSys"], float64(rtm.MSpanSys))
	if err != nil {
		return
	}

	err = updateGaugeMetric(m.metrics["NextGC"], float64(rtm.NextGC))
	if err != nil {
		return
	}

	err = updateGaugeMetric(m.metrics["NumForcedGC"], float64(rtm.NumForcedGC))
	if err != nil {
		return
	}
	err = updateGaugeMetric(m.metrics["NumGC"], float64(rtm.NumGC))
	if err != nil {
		return
	}
	err = updateGaugeMetric(m.metrics["OtherSys"], float64(rtm.OtherSys))
	if err != nil {
		return
	}
	err = updateGaugeMetric(m.metrics["PauseTotalNs"], float64(rtm.PauseTotalNs))
	if err != nil {
		return
	}

	err = updateGaugeMetric(m.metrics["StackInuse"], float64(rtm.StackInuse))
	if err != nil {
		return
	}

	err = updateGaugeMetric(m.metrics["StackSys"], float64(rtm.StackSys))
	if err != nil {
		return
	}

	err = updateGaugeMetric(m.metrics["Sys"], float64(rtm.Sys))
	if err != nil {
		return
	}

	err = updateGaugeMetric(m.metrics["Mallocs"], float64(rtm.Mallocs))
	if err != nil {
		return
	}

	err = updateGaugeMetric(m.metrics["TotalAlloc"], float64(rtm.TotalAlloc))
	if err != nil {
		return
	}

	err = updateCounterMetric(m.metrics["PollCount"], 1)
	if err != nil {
		return
	}

	err = updateGaugeMetric(m.metrics["RandomValue"], rand.Float64())
	if err != nil {
		return
	}

	return
}

func (m *monitor) PollAdditionalMetrics() (err error) {
	fmt.Println("poll additional metrics")
	v, err := mem.VirtualMemory()

	if err != nil {
		return fmt.Errorf("failed to get memory metrics %w", err)
	}

	err = updateGaugeMetric(m.metrics["TotalMemory"], float64(v.Total))
	if err != nil {
		return fmt.Errorf("failed to update TotalMemory metric %w", err)
	}

	err = updateGaugeMetric(m.metrics["FreeMemory"], float64(v.Free))
	if err != nil {
		return fmt.Errorf("failed to update FreeMemory metric %w", err)
	}

	cpuStat, err := cpu.Times(false)

	if err != nil {
		return fmt.Errorf("failed to get cpu metrics %w", err)
	}

	err = updateGaugeMetric(m.metrics["CPUutilization1"], cpuStat[0].User)

	if err != nil {
		return fmt.Errorf("failed to update CPUutilization1 metric %w", err)
	}

	return
}

func (m *monitor) ReportMetrics() (err error) {
	fmt.Println("reportMetrics")

	m.wg.Add(1)
	defer m.wg.Done()

	if len(m.metrics) == 0 {
		fmt.Println("empty batch")
		return
	}

	var metricsBatch []metric.Metrics

	for _, value := range m.metrics {
		metricsBatch = append(metricsBatch, *value)
	}

	metricsBytes, err := json.Marshal(metricsBatch)
	if err != nil {
		return fmt.Errorf("failed to marshal metrics batch: %w", err)
	}

	if m.cryptoKey != "" {
		dat, err := os.ReadFile(m.cryptoKey)
		if err != nil {
			return fmt.Errorf("failed to read public key from file: %w", err)
		}

		block, _ := pem.Decode(dat)

		publicKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
		if err != nil {
			return fmt.Errorf("failed to parse public key: %w", err)
		}

		//The message must be no longer than the length of the public modulus minus 11 bytes.

		step := publicKey.Size() - 11

		var encryptedBytes []byte

		for start := 0; start < len(metricsBytes); start += step {
			finish := start + step
			if finish > len(metricsBytes) {
				finish = len(metricsBytes)
			}

			encryptedBlockBytes, err := rsa.EncryptPKCS1v15(cryptoRand.Reader, publicKey, metricsBytes[start:finish])

			if err != nil {
				return fmt.Errorf("failed to encrypt request body: %w", err)
			}

			encryptedBytes = append(encryptedBytes, encryptedBlockBytes...)
		}

		metricsBytes = encryptedBytes
	}

	client := resty.New()
	client.SetPreRequestHook(func(c *resty.Client, req *http.Request) (err error) {

		fmt.Println("On before request")

		if m.key != "" {
			hash := sha256.New()

			bodyBytes, errBody := io.ReadAll(req.Body)

			if errBody != nil {
				return fmt.Errorf("failed to read request body: %w", errBody)
			}

			req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

			keyBytes := []byte(m.key)

			src := append(bodyBytes, keyBytes...)
			hash.Write(src)

			dst := hash.Sum(nil)

			encodedDst := base64.StdEncoding.EncodeToString(dst)

			req.Header.Set("HashSHA256", encodedDst)

			fmt.Println("HashSHA256 set to ", encodedDst)
		}

		return
	})

	realIP, err := getIP()
	if err != nil {
		return fmt.Errorf("failed to get IP: %w", err)
	}

	_, err = client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept-Encoding", "gzip").
		SetHeader("X-Real-Ip", realIP.IP.String()).
		SetHeader("Rsa-Encrypted", strconv.FormatBool(m.cryptoKey != "")).
		SetBody(metricsBytes).
		Post(fmt.Sprintf("%v/updates/", m.hostAddress))

	//TODO: properly handle connection refused error (don't quit goroutine)
	//but quit on fatal error

	if err != nil {
		return fmt.Errorf("failed to post updates: %w", err)
	}
	return
}

// ReportMetricsGrpc sends metrics by gRPC
func (m *monitor) ReportMetricsGrpc() (err error) {
	fmt.Println("reportMetrics gRPC")

	if len(m.metrics) == 0 {
		fmt.Println("empty batch")
		return
	}

	var opts []grpc.DialOption

	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.Dial(m.hostGrpc, opts...)
	if err != nil {
		return fmt.Errorf("failed to dial gRPC server: %w", err)
	}

	defer conn.Close()

	m.wg.Add(1)
	defer m.wg.Done()

	var metricsBatch pb.MetricsBatchRequest

	for _, value := range m.metrics {
		pbMetric := &pb.Metric{
			Id:    value.ID,
			Type:  value.MType,
			Delta: value.Delta,
			Value: value.Value,
		}
		metricsBatch.Metrics = append(metricsBatch.Metrics, pbMetric)
	}

	metricsBatch.Count = int32(len(metricsBatch.Metrics))

	c := pb.NewMetricsServiceV1Client(conn)

	ctx := context.TODO()
	resp, err := c.MetricsBatchV1(ctx, &metricsBatch)

	if err != nil {
		return fmt.Errorf("failed to post updates: %w", err)
	}

	fmt.Println(resp)
	return
}

func getIP() (realIP *net.IPNet, err error) {

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		err = fmt.Errorf("failed to get IPs: %w", err)
		return
	}

	for _, addr := range addrs {
		ipNet, ok := addr.(*net.IPNet)
		if ok && !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
			realIP = ipNet
			break
		}
	}
	return
}
