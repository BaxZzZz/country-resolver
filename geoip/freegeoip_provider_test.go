package geoip

import (
	"testing"
	"net/http"
	"net/http/httptest"
)

const freeGeoIpJson =
`{
"ip":"192.30.253.112",
"country_code":"US",
"country_name":"United States",
"region_code":"CA",
"region_name":"California",
"city":"San Francisco",
"zip_code":"94107",
"time_zone":"America/Los_Angeles",
"latitude":37.7697,
"longitude":-122.3933,
"metro_code":807
}`

type FakeClient struct {
	Result []byte
	Error error
}

func (client* FakeClient)Request(url string) ([]byte, error) {
	if client.Error != nil {
		return nil, client.Error
	}

	return client.Result, nil
}

func TestFreeGeoIPRequest(t *testing.T) {
	handler := func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte(freeGeoIpJson))
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	fakeClient := FakeClient{ Result: []byte(freeGeoIpJson) }
	provider := freeGeoIPProvider{ client: &fakeClient }

	info, err := provider.GetIpInfo("123")

	if err != nil {
		t.Fatalf("GetIpInfo: %v", err)
	}

	if info.CountryName != "United States" {
		t.Fatal("CountryName: ",  info.CountryName)
	}
}