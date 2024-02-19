package main

import (
	"github.com/lavandosovich/goya/internal"
	"log"
	"net/http"
)

func main() {
	memStorage := internal.NewMemStorage()

	http.Handle("/update/", internal.HandlerWrapper(memStorage, internal.PostHandler))

	server := &http.Server{
		Addr:    "127.0.0.1:8080",
		Handler: nil,
	}

	// вызов ListenAndServe — блокирующий, последний в программе
	// возникающие ошибки на серверных машинах пишут в системный лог,
	// а не в стандартную консоль ошибок,
	// поэтому обычно вызывают вот так
	log.Fatal(server.ListenAndServe())
}
