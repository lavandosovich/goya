package internal

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"slices"
	"strconv"
	"strings"
)

func MetricTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		metricsTypes := []string{"counter", "gauge"}
		metricType := chi.URLParam(r, string(MetricType))

		if !slices.Contains(metricsTypes, metricType) {
			w.WriteHeader(http.StatusNotImplemented)
			w.Write([]byte("fail"))
			return
		}
		ctx := context.WithValue(r.Context(), MetricType, metricType)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RootHandler(w http.ResponseWriter, _ *http.Request, storage *MemStorage) {
	w.Header().Set("content-type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(storage.ReduceMetricsToHTML().Bytes())
}

func GetHandler(w http.ResponseWriter, r *http.Request, storage *MemStorage) {
	var (
		metricValue   string
		ok            bool
		counterMetric Counter
		gaugeMetric   Gauge
	)
	metricType := r.
		Context().
		Value(MetricType)
	metricName := chi.URLParam(r, string(MetricName))

	if metricType == "counter" {
		counterMetric, ok = storage.GetCounterMetric(metricName)
		metricValue = strconv.Itoa(int(counterMetric))
	} else {
		gaugeMetric, ok = storage.GetGaugeMetric(metricName)
		metricValue =
			strings.TrimRight(fmt.Sprintf("%f", float64(gaugeMetric)), "0")
	}
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(""))
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(metricValue))
}

func PostHandler(w http.ResponseWriter, r *http.Request, storage *MemStorage) {
	ctx := r.Context()
	metricType := ctx.Value(MetricType)
	metricValue := chi.URLParam(r, string(MetricValue))
	metricName := chi.URLParam(r, string(MetricName))

	w.Header().Set("content-type", "application/text")

	if metricType == "counter" {
		value, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("wrong type"))
			return
		}
		storage.SetCounterMetric(metricName, Counter(value))
	} else {
		s, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("wrong type"))
			return
		}
		storage.SetGaugeMetric(metricName, Gauge(s))
	}
	log := fmt.Sprintf("%s %d\n", r.URL.Path, http.StatusAccepted)
	fmt.Print(log)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}
