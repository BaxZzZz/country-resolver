package geoip

import (
	"github.com/golang/mock/gomock"
	"testing"
	"time"
)

func TestRequestGetIPInfo(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	ipAddr := "8.8.8.8"
	expectedIPInfo := &IPInfo{CountryName: "Test"}

	provider := NewMockProvider(mockCtrl)
	provider.EXPECT().GetIPInfo(ipAddr).Return(expectedIPInfo, nil)

	request, err := NewRequest([]Provider{provider}, 10, time.Duration(1000)*time.Minute)

	if err != nil {
		t.Fatal(err)
	}

	info, err := request.GetIPInfo(ipAddr)

	if err != nil {
		t.Fatal(err)
	}

	if info.CountryName != expectedIPInfo.CountryName {
		t.Fatalf("IP info not equal, expected: %v, actual: %v", *expectedIPInfo, *info)
	}
}

func TestNextProviderWhenRequestsLimit(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	ipAddr := "8.8.8.8"
	expectedFirstIPInfo := &IPInfo{CountryName: "Test1"}
	expectedSecondIPInfo := &IPInfo{CountryName: "Test2"}

	firstProvider := NewMockProvider(mockCtrl)
	firstProvider.EXPECT().GetIPInfo(ipAddr).Return(expectedFirstIPInfo, nil)

	secondProvider := NewMockProvider(mockCtrl)
	secondProvider.EXPECT().GetIPInfo(ipAddr).Return(expectedSecondIPInfo, nil)

	request, err := NewRequest([]Provider{firstProvider, secondProvider}, 1,
		time.Duration(1000)*time.Minute)
	if err != nil {
		t.Fatal(err)
	}

	firstInfo, err := request.GetIPInfo(ipAddr)
	if err != nil {
		t.Fatal(err)
	}

	if firstInfo.CountryName != expectedFirstIPInfo.CountryName {
		t.Fatalf("IP info not equal, expected: %v, actual: %v", *expectedFirstIPInfo, *firstInfo)
	}

	secondInfo, err := request.GetIPInfo(ipAddr)
	if err != nil {
		t.Fatal(err)
	}

	if secondInfo.CountryName != expectedSecondIPInfo.CountryName {
		t.Fatalf("IP info not equal, expected: %v, actual: %v", *expectedSecondIPInfo, *firstInfo)
	}
}

func TestWhenRequestsLessWhenLimits(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	ipAddr := "8.8.8.8"
	expectedFirstIPInfo := &IPInfo{CountryName: "Test1"}
	expectedSecondIPInfo := &IPInfo{CountryName: "Test2"}

	firstProvider := NewMockProvider(mockCtrl)
	secondProvider := NewMockProvider(mockCtrl)

	firstProvider.EXPECT().GetIPInfo(ipAddr).Return(expectedFirstIPInfo, nil)
	firstProvider.EXPECT().GetIPInfo(ipAddr).Return(expectedSecondIPInfo, nil)

	request, err := NewRequest([]Provider{firstProvider, secondProvider}, 2,
		time.Duration(1)*time.Millisecond)
	if err != nil {
		t.Fatal(err)
	}

	firstInfo, err := request.GetIPInfo(ipAddr)
	if err != nil {
		t.Fatal(err)
	}

	if firstInfo.CountryName != expectedFirstIPInfo.CountryName {
		t.Fatalf("IP info not equal, expected: %v, actual: %v", *expectedFirstIPInfo, *firstInfo)
	}

	secondInfo, err := request.GetIPInfo(ipAddr)
	if err != nil {
		t.Fatal(err)
	}

	if secondInfo.CountryName != expectedSecondIPInfo.CountryName {
		t.Fatalf("IP info not equal, expected: %v, actual: %v", *expectedSecondIPInfo, *firstInfo)
	}
}

func TestWhenAllProvidersUnavailable(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	ipAddr := "8.8.8.8"

	firstProvider := NewMockProvider(mockCtrl)
	firstProvider.EXPECT().GetIPInfo(gomock.Any()).Return(&IPInfo{}, nil).AnyTimes()

	secondProvider := NewMockProvider(mockCtrl)
	secondProvider.EXPECT().GetIPInfo(gomock.Any()).Return(&IPInfo{}, nil).AnyTimes()

	request, err := NewRequest([]Provider{firstProvider, secondProvider}, 2,
		time.Duration(500) * time.Millisecond)
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 4; i++ {
		_, err := request.GetIPInfo(ipAddr)
		if err != nil {
			t.Fatal(err)
		}
	}

	_, err = request.GetIPInfo(ipAddr)
	if err == nil {
		t.Fatal("Providers must be unavailable")
	}

	time.Sleep(1*time.Second)

	isAvailable := false
	for i := 0; i < 1000000; i++ {
		_, err := request.GetIPInfo(ipAddr)
		if err == nil {
			isAvailable = true
		}
	}

	if !isAvailable {
		t.Fatal("Providers must be available")
	}
}
