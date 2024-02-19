package internal

import (
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"
)

func PostHandler(w http.ResponseWriter, r *http.Request, storage *MemStorage) {
	switch r.Method {
	case http.MethodPost:
		fmt.Println(r.URL.Path)
		metricsTypes := []string{"counter", "gauge"}
		splittedPath := strings.Split(r.URL.Path, "/")
		w.Header().Set("content-type", "application/text")

		if len(splittedPath) != 5 {
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Printf(fmt.Sprintf("%s %d\n", r.URL.Path, http.StatusMethodNotAllowed))
			w.Write([]byte("fail"))
			return
		}

		if !slices.Contains(metricsTypes, splittedPath[2]) {
			w.WriteHeader(http.StatusNotImplemented)
			w.Write([]byte("fail"))
			return
		}

		if splittedPath[2] == "counter" {
			value, err := strconv.ParseInt(splittedPath[4], 10, 64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("wrong type"))
				return
			}
			(*storage).SetMetric(splittedPath[3], Counter(value))
		} else {
			s, err := strconv.ParseFloat(splittedPath[4], 64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("wrong type"))
				return
			}
			(*storage).SetMetric(splittedPath[3], Gauge(s))
		}
		fmt.Printf(fmt.Sprintf("%s %d\n", r.URL.Path, http.StatusAccepted))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
