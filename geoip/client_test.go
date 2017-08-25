package geoip

import (
	"testing"
	"net/http"
	"net/http/httptest"
)

func TestHttpClientRequest(t *testing.T) {
	handler := func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("OK"))
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	client := httpClient{}
	resp, err := client.Request(server.URL)
	if err != nil {
		t.Fatalf("Request: %v", err)
	}

	if string(resp) != "OK" {
		t.Fatal("Body: ",  string(resp))
	}
}