package cache

import (
	"time"
)

type Item struct {
	Key       string
	Object    interface{}
	Timestamp int64
}

type DataStore interface {
	GetItemsCount() (uint, error)
	SetItem(item *Item) error
	FindItem(key string) (*Item, error)
	RemoveOldItem() error
	Close()
}

func makeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
