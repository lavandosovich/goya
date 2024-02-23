package main

import (
	"time"

	"github.com/lavandosovich/goya/internal"
)

func main() {
	internal.PollMetrics(time.Second*2, time.Second*10, internal.POSTMetrics, "http://localhost:8080")
}
