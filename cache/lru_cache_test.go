package cache

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestLRUCacheGetNonExisten(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	expectedError := errors.New("Test")

	store := NewMockDataStore(mockCtrl)
	store.EXPECT().Find("123", gomock.Any()).Return(expectedError)

	lruCache := LRUCache{
		ItemsLimit: 10,
		Store:      store,
	}

	err := lruCache.Get("123", "")
	if err != expectedError {
		t.Fatalf("Must be error, expected: %v, actual: %v", expectedError, err)
	}
}

func TestLRUCacheGet(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	data := "1234"
	store := NewMockDataStore(mockCtrl)
	store.EXPECT().Find("123", data).Return(nil)

	lruCache := LRUCache{
		ItemsLimit: 10,
		Store:      store,
	}

	err := lruCache.Get("123", data)

	if err != nil {
		t.Fatal(err)
	}
}

func TestLRUCacheSet(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	store := NewMockDataStore(mockCtrl)
	store.EXPECT().GetCount().Return(uint(1), nil)
	store.EXPECT().Set("123", "1234").Return(nil)

	lruCache := LRUCache{
		ItemsLimit: 10,
		Store:      store,
	}

	err := lruCache.Set("123", "1234")

	if err != nil {
		t.Fatal(err)
	}
}

func TestLRUCacheWhenLimitExceeded(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	store := NewMockDataStore(mockCtrl)
	store.EXPECT().GetCount().Return(uint(10), nil)
	store.EXPECT().RemoveOld().Return(nil)
	store.EXPECT().Set("123", "1234").Return(nil)

	lruCache := LRUCache{
		ItemsLimit: 10,
		Store:      store,
	}

	err := lruCache.Set("123", "1234")

	if err != nil {
		t.Fatal(err)
	}
}
