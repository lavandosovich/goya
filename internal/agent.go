package internal

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"
)

func POSTMetrics(metrics *Metrics, address string) ([]*http.Response, []error) {
	var requests []*http.Request
	var requestErrors []error
	var responses []*http.Response

	client := &http.Client{}

	request, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("%s/update/counter/pollcount/%d", address, metrics.PollCount),
		nil)
	if err != nil {
		return nil, append(requestErrors, err)
	}
	request.Header.Add("Content-Type", "text/plain")

	requests = append(requests, request)

	v := reflect.ValueOf(*metrics)
	typeOfS := v.Type()

	for i := 0; i < v.NumField(); i++ {
		if typeOfS.Field(i).Name == "PollCounter" {
			continue
		}
		request, err := http.NewRequest(
			http.MethodPost,
			fmt.Sprintf(
				"%s/update/gauge/%s/%v",
				address,
				strings.ToLower(typeOfS.Field(i).Name),
				v.Field(i).Interface(),
			),
			nil)
		if err != nil {
			return nil, append(requestErrors, err)
		}
		request.Header.Add("Content-Type", "text/plain")
		requests = append(requests, request)
	}

	for _, request := range requests {
		go func(request *http.Request) {
			response, err := client.Do(request)
			if err != nil {
				requestErrors = append(requestErrors, err)
			}
			defer response.Body.Close()
			responses = append(responses, response)
		}(request)
	}

	return responses, requestErrors
}
