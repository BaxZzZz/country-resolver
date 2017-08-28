package geoip

import (
	"github.com/golang/mock/gomock"
	"testing"
)

const freeGeoIPJson = `{
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

func TestFreeGeoIPGetIPInfo(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	expectInfo := &IPInfo{
		CountryName: "United States",
	}

	mockClient := NewMockClient(mockCtrl)
	mockClient.EXPECT().Request(gomock.Any()).Return([]byte(freeGeoIPJson), nil)

	provider := freeGeoIPProvider{client: mockClient}

	info, err := provider.GetIPInfo("8.8.8.8")
	if err != nil {
		t.Fatalf("GetIpInfo: %v", err)
	}

	if *expectInfo != *info {
		t.Fatalf("Not equal expect: %v, actual: %v", *expectInfo, *info)
	}
}
