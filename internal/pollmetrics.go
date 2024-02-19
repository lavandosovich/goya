package internal

import (
	"fmt"
	"math/rand"
	"runtime"
	"time"
)

func getMetrics(pollCount Counter) *Metrics {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	metrics := &Metrics{
		Alloc:         Gauge(float64(m.Alloc)),
		BuckHashSys:   Gauge(float64(m.BuckHashSys)),
		Frees:         Gauge(float64(m.Frees)),
		GCCPUFraction: Gauge(m.GCCPUFraction),
		GCSys:         Gauge(float64(m.GCSys)),
		HeapAlloc:     Gauge(float64(m.HeapAlloc)),
		HeapIdle:      Gauge(float64(m.HeapIdle)),
		HeapInuse:     Gauge(float64(m.HeapInuse)),
		HeapObjects:   Gauge(float64(m.HeapObjects)),
		HeapReleased:  Gauge(float64(m.HeapReleased)),
		HeapSys:       Gauge(float64(m.HeapSys)),
		LastGC:        Gauge(float64(m.LastGC)),
		Lookups:       Gauge(float64(m.Lookups)),
		MCacheInuse:   Gauge(float64(m.MCacheInuse)),
		MCacheSys:     Gauge(float64(m.MCacheSys)),
		MSpanInuse:    Gauge(float64(m.MSpanInuse)),
		MSpanSys:      Gauge(float64(m.MSpanSys)),
		Mallocs:       Gauge(float64(m.Mallocs)),
		NextGC:        Gauge(float64(m.NextGC)),
		NumForcedGC:   Gauge(float64(m.NumForcedGC)),
		NumGC:         Gauge(float64(m.NumGC)),
		OtherSys:      Gauge(float64(m.OtherSys)),

		PauseTotalNs: Gauge(float64(m.PauseTotalNs)),
		StackInuse:   Gauge(float64(m.StackInuse)),
		StackSys:     Gauge(float64(m.StackSys)),
		Sys:          Gauge(float64(m.Sys)),
		TotalAlloc:   Gauge(float64(m.TotalAlloc)),
		PollCount:    pollCount,
		RandomValue:  Gauge(rand.Int63n(1_000_000_000)),
	}

	return metrics
}

func PollMetrics(pollDuration, reportDuration time.Duration, reporterFunc MetricsReporter, address string) {
	pollTicker := time.NewTicker(pollDuration)
	reportTicker := time.NewTicker(reportDuration)
	startTime := time.Now()
	var pollCount Counter = 0
	var neededMetrics *Metrics

	go func() {
		for {
			tick := <-reportTicker.C
			if neededMetrics == nil {
				continue
			}
			responses, errors := reporterFunc(neededMetrics, address)

			if errors != nil {
				fmt.Println(errors)
				panic("error on reporting metrics")
			}

			for _, response := range responses {
				_ = response.Body.Close()
			}

			fmt.Println("From goroutine", int(tick.Sub(startTime).Seconds()))
		}
	}()

	for i := 0; ; i++ {
		tick := <-pollTicker.C
		pollCount += 1
		neededMetrics = getMetrics(pollCount)
		fmt.Println(int(tick.Sub(startTime).Seconds()))
	}
}
