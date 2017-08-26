package geoip

import (
	"reflect"
	"testing"
)

func TestNewProviders(t *testing.T) {
	providers, err := NewProviders([]string{FREE_GEO_IP_NAME, NEKUDO_NAME})

	if err != nil {
		t.Fatalf("Failed to create providers, error: %v", err)
	}

	if len(providers) != 2 {
		t.Fatal("Failed providers count")
	}

	actualType := reflect.TypeOf(providers[0]).Elem().Name()
	expectedType := reflect.TypeOf(freeGeoIPProvider{}).Name()

	if actualType != expectedType {
		t.Fatal("Invalid type, expected: " + expectedType + ", actual: " + actualType)
	}

	actualType = reflect.TypeOf(providers[1]).Elem().Name()
	expectedType = reflect.TypeOf(nekudoProvider{}).Name()

	if actualType != expectedType {
		t.Fatal("Invalid type, expected: " + expectedType + ", actual: " + actualType)
	}
}
