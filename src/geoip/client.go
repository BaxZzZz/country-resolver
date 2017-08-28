package geoip

import (
	"io/ioutil"
	"net/http"
)

// Client interface
type Client interface {
	Request(url string) ([]byte, error)
}

// HTTP client implementation
type httpClient struct {
	client http.Client
}

// Request data through HTTP protocol
func (client *httpClient) Request(url string) ([]byte, error) {
	resp, err := client.client.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
