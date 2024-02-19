package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/lavandosovich/goya/internal"
)

func main() {
	memStorage := internal.NewMemStorage()

	http.Handle("/update/counter/PollCount", internal.HandlerWrapper(memStorage, internal.PostHandler))
	for _, metricType := range internal.GetMetricNames() {
		metrics := fmt.Sprintf("/update/gauge/%s/", strings.ToLower(metricType))
		http.Handle(metrics, internal.HandlerWrapper(memStorage, internal.PostHandler))
	}
	server := &http.Server{
		Addr: "127.0.0.1:8080",
	}
	// вызов ListenAndServe — блокирующий, последний в программе
	// возникающие ошибки на серверных машинах пишут в системный лог,
	// а не в стандартную консоль ошибок,
	// поэтому обычно вызывают вот так
	log.Fatal(server.ListenAndServe())
}
