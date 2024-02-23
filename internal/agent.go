package internal

import (
	"fmt"
	"net/http"
	"reflect"
	"sync"
)

func createRequest(method, url, contentType string) (*http.Request, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", contentType)
	return req, nil
}

func sendRequest(client *http.Client, request *http.Request) error {
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	return nil
}

func POSTMetrics(metrics Metrics, address string) []error {
	var requests []*http.Request
	var requestErrors []error

	client := &http.Client{}

	request, err := createRequest(
		http.MethodPost,
		fmt.Sprintf("%s/update/counter/pollcount/%d", address, metrics.PollCount),
		"text/plain",
	)

	if err != nil {
		return append(requestErrors, err)
	}
	requests = append(requests, request)

	v := reflect.ValueOf(metrics)
	typeOfS := v.Type()

	for i := 0; i < v.NumField(); i++ {
		if typeOfS.Field(i).Name == "PollCounter" {
			continue
		}
		request, err := createRequest(
			http.MethodPost,
			fmt.Sprintf(
				"%s/update/gauge/%s/%v",
				address,
				typeOfS.Field(i).Name,
				v.Field(i).Interface(),
			),
			"text/plain",
		)
		if err != nil {
			return append(requestErrors, err)
		}
		request.Header.Add("Content-Type", "text/plain")
		requests = append(requests, request)
	}

	var wg sync.WaitGroup
	errorChan := make(chan error, len(requests))
	for _, request := range requests {
		wg.Add(1)
		go func(request *http.Request) {
			errorChan <- sendRequest(client, request)
			defer wg.Done()
		}(request)
	}
	wg.Wait()
	close(errorChan)

	for err := range errorChan {
		if err != nil {
			requestErrors = append(requestErrors, err)
		}
	}
	return requestErrors
}
