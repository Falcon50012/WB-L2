package wgetm

import (
	"fmt"
	"io"
	"net/http"
)

const userAgent = "GoWget/1.0"

type fetchResult struct {
	body        []byte
	contentType string
	statusCode  int
}

func fetch(client *http.Client, rawURL string) (fetchResult, error) {
	req, err := http.NewRequest(http.MethodGet, rawURL, nil)
	if err != nil {
		return fetchResult{}, fmt.Errorf("request error: %w", err)
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,*/*;q=0.8")

	resp, err := client.Do(req)
	if err != nil {
		return fetchResult{}, fmt.Errorf("GET %s: %w", rawURL, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 50<<20))
	if err != nil {
		return fetchResult{}, fmt.Errorf("body read error: %w", err)
	}

	return fetchResult{
		body:        body,
		contentType: resp.Header.Get("Content-Type"),
		statusCode:  resp.StatusCode,
	}, nil
}
