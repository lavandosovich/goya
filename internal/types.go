package internal

import (
	"bytes"
	"fmt"
	"net/http"
)

type metricKey string

const (
	MetricType  metricKey = "metricType"
	MetricValue metricKey = "metricValue"
	MetricName  metricKey = "metricName"
)

type Gauge float64
type Counter int64

type MetricI interface {
	Gauge | Counter
}
type Metrics struct {
	Alloc         Gauge
	BuckHashSys   Gauge
	Frees         Gauge
	GCCPUFraction Gauge
	GCSys         Gauge
	HeapAlloc     Gauge
	HeapIdle      Gauge
	HeapInuse     Gauge
	HeapObjects   Gauge
	HeapReleased  Gauge
	HeapSys       Gauge
	LastGC        Gauge
	Lookups       Gauge
	MCacheInuse   Gauge
	MCacheSys     Gauge
	MSpanInuse    Gauge
	MSpanSys      Gauge
	Mallocs       Gauge
	NextGC        Gauge
	NumForcedGC   Gauge
	NumGC         Gauge
	OtherSys      Gauge
	PauseTotalNs  Gauge
	StackInuse    Gauge
	StackSys      Gauge
	Sys           Gauge
	TotalAlloc    Gauge
	PollCount     Counter
	RandomValue   Gauge
}

type MemStorage struct {
	counterStorage map[string]Counter
	gaugeStorage   map[string]Gauge
}

type IMemStorage interface {
	SetMetric(metricName string, metricValue any)
	GetMetric(metricName string) any
	SetCounterMetric(metricName string, metricValue Counter)
	GetCounterMetric(metricName string) Counter
	SetGaugeMetric(metricName string, metricValue Gauge)
	GetGaugeMetric(metricName string) Gauge
	ReduceMetricsToHtml() string
}

type MetricsReporter func(metrics *Metrics, address string) []error

func NewMemStorage() *MemStorage {
	return &MemStorage{
		counterStorage: make(map[string]Counter),
		gaugeStorage:   make(map[string]Gauge),
	}
}

func (memStorage MemStorage) SetCounterMetric(metricName string, metricValue Counter) {
	memStorage.counterStorage[metricName] = metricValue
}

func (memStorage MemStorage) GetCounterMetric(metricName string) Counter {
	return memStorage.counterStorage[metricName]
}

func (memStorage MemStorage) SetGaugeMetric(metricName string, metricValue Gauge) {
	memStorage.gaugeStorage[metricName] = metricValue
}

func (memStorage MemStorage) GetGaugeMetric(metricName string) Gauge {
	return memStorage.gaugeStorage[metricName]
}

func (memStorage MemStorage) ReduceMetricsToHTML() *bytes.Buffer {
	var htmlBody bytes.Buffer
	for k, v := range memStorage.gaugeStorage {
		htmlBody.WriteString(fmt.Sprintf("<div>%s: %f</div>\n", k, v))
	}
	for k, v := range memStorage.counterStorage {
		htmlBody.WriteString(fmt.Sprintf("<div>%s: %d</div>\n", k, v))
	}
	return &htmlBody
}

func HandlerWrapper(
	storage *MemStorage,
	function func(http.ResponseWriter, *http.Request, *MemStorage)) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		function(w, r, storage)
	}
}
