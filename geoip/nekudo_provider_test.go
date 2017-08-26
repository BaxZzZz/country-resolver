package geoip

import (
	"testing"

	"github.com/golang/mock/gomock"
)

const nekudoIPJson = `{
   "city":false,
   "country":{
      "name":"United States",
      "code":"US"
   },
   "location":{
      "accuracy_radius":1000,
      "latitude":37.751,
      "longitude":-97.822
   },
   "ip":"8.8.8.8"
}`

func TestNekudoGetIPInfo(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	expectInfo := &IPInfo{
		CountryName: "United States",
	}

	mockClient := NewMockClient(mockCtrl)
	mockClient.EXPECT().Request(gomock.Any()).Return([]byte(nekudoIPJson), nil)

	provider := nekudoProvider{client: mockClient}

	info, err := provider.GetIpInfo("8.8.8.8")
	if err != nil {
		t.Fatalf("GetIpInfo: %v", err)
	}

	if *expectInfo != *info {
		t.Fatalf("Not equal expect: %v, actual: %v", *expectInfo, *info)
	}
}
