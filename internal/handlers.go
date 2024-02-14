package internal

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func PostHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		splittedPath := strings.Split(r.URL.Path, "/")
		w.Header().Set("content-type", "application/text")

		if len(splittedPath) != 5 {
			w.WriteHeader(404)
			fmt.Println(fmt.Sprintf("%s %d", r.URL.Path, http.StatusNotFound))

			w.Write([]byte("not 4"))
			return
		}

		_, err := strconv.Atoi(splittedPath[4])
		if err != nil {
			fmt.Println(fmt.Sprintf("%s %d", r.URL.Path, http.StatusNotFound))
			w.Write([]byte("not int"))

			return
		}
		fmt.Println(fmt.Sprintf("%s %d", r.URL.Path, http.StatusAccepted))
		w.WriteHeader(200)
	default:
		w.WriteHeader(405)
	}
}
