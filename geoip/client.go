package geoip

import (
	"io/ioutil"
	"net/http"
)

type Client interface {
	Request(url string) ([]byte, error)
}

type httpClient struct {
	client http.Client
}

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
