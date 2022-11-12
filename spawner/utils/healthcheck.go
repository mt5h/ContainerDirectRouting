package utils

import (
	"context"
	"log"
	"net/http"
	"time"
)

func HttpHealthCheck(endpointUrl string, timeout time.Duration, acceptedStatusCode int) bool {

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	log.Println("Making request to", endpointUrl)

	req, _ := http.NewRequestWithContext(ctx, "GET", endpointUrl, nil)

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		return false
	}

	log.Println("Got a response with status", resp.Status)

	if resp.StatusCode == acceptedStatusCode {
		return true
	}

	return false
}
