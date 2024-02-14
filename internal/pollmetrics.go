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
		Alloc:         Gauge(int(m.Alloc)),
		BuckHashSys:   Gauge(int(m.BuckHashSys)),
		Frees:         Gauge(int(m.Frees)),
		GCCPUFraction: Gauge(int(m.GCCPUFraction)),
		GCSys:         Gauge(int(m.GCSys)),
		HeapAlloc:     Gauge(int(m.HeapAlloc)),
		HeapIdle:      Gauge(int(m.HeapIdle)),
		HeapInuse:     Gauge(int(m.HeapInuse)),
		HeapObjects:   Gauge(int(m.HeapObjects)),
		HeapReleased:  Gauge(int(m.HeapReleased)),
		HeapSys:       Gauge(int(m.HeapSys)),
		LastGC:        Gauge(int(m.LastGC)),
		Lookups:       Gauge(int(m.Lookups)),
		MCacheInuse:   Gauge(int(m.MCacheInuse)),
		MCacheSys:     Gauge(int(m.MCacheSys)),
		MSpanInuse:    Gauge(int(m.MSpanInuse)),
		MSpanSys:      Gauge(int(m.MSpanSys)),
		Mallocs:       Gauge(int(m.Mallocs)),
		NextGC:        Gauge(int(m.NextGC)),
		NumForcedGC:   Gauge(int(m.NumForcedGC)),
		NumGC:         Gauge(int(m.NumGC)),
		OtherSys:      Gauge(int(m.OtherSys)),

		PauseTotalNs: Gauge(int(m.PauseTotalNs)),
		StackInuse:   Gauge(int(m.StackInuse)),
		StackSys:     Gauge(int(m.StackSys)),
		Sys:          Gauge(int(m.Sys)),
		TotalAlloc:   Gauge(int(m.TotalAlloc)),
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
			fmt.Println(responses)
			fmt.Println("From goroutine", int(tick.Sub(startTime).Seconds()))
		}
	}()

	for i := 0; ; i++ {
		tick := <-pollTicker.C
		pollCount += 1
		neededMetrics = getMetrics(pollCount)
		fmt.Println("From main routine")
		fmt.Println(neededMetrics)
		fmt.Println(int(tick.Sub(startTime).Seconds()))
	}
}
