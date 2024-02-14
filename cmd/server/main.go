package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/lavandosovich/goya/internal"
)

func main() {
	http.Handle("/update/counter/pollcount", http.HandlerFunc(internal.PostHandler))
	for _, metricType := range internal.GetMetricNames() {
		metrics := fmt.Sprintf("/update/gauge/%s/", strings.ToLower(metricType))
		http.Handle(
			metrics,
			http.HandlerFunc(internal.PostHandler))

	}
	server := &http.Server{
		Addr: "127.0.0.1:8081",
	}
	// вызов ListenAndServe — блокирующий, последний в программе
	// возникающие ошибки на серверных машинах пишут в системный лог,
	// а не в стандартную консоль ошибок,
	// поэтому обычно вызывают вот так
	log.Fatal(server.ListenAndServe())
}
