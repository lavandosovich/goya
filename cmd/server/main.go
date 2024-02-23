package main

import (
	"github.com/lavandosovich/goya/internal"
	"log"
	"net/http"
)

func main() {
	memStorage := internal.NewMemStorage()

	server := &http.Server{
		Addr:    "127.0.0.1:8080",
		Handler: internal.InitChiRouter(memStorage),
	}

	log.Fatal(server.ListenAndServe())
}
