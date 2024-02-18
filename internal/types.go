package internal

import (
	"net/http"
	"reflect"
)

type Gauge float64
type Counter int64

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
	storage map[string]any
}

type IMemStorage interface {
	SetMetric(metricName string, metricValue any)
	GetMetric(metricName string) any
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		storage: make(map[string]any),
	}
}

func (memStorage MemStorage) SetMetric(metricName string, metricValue any) {
	memStorage.storage[metricName] = metricValue
}

func (memStorage MemStorage) GetMetric(metricName string) any {
	return memStorage.storage[metricName]
}

func HandlerWrapper(
	storage *MemStorage,
	function func(http.ResponseWriter, *http.Request, *MemStorage)) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		function(w, r, storage)
	}
}

type MetricsReporter func(metrics *Metrics, address string) ([]*http.Response, []error)

func GetMetricNames() []string {
	var fields []string
	val := reflect.Indirect(reflect.ValueOf(Metrics{})).Type()

	for i := 0; i < val.NumField(); i++ {
		name := val.Field(i).Name
		if name == "PollCount" {
			continue
		}
		fields = append(fields, name)
	}

	return fields
}
