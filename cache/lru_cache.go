package cache

import "errors"

type LRUCache struct {
	ItemsLimit uint
	Store      DataStore
}

func (lruCache *LRUCache) Set(key string, object interface{}) error {
	countItems, err := lruCache.Store.GetCount()
	if err != nil {
		return err
	}

	if countItems >= lruCache.ItemsLimit {
		err = lruCache.Store.RemoveOld()
		if err != nil {
			return err
		}
	}

	return lruCache.Store.Set(key, object)
}

func (lruCache *LRUCache) Get(key string, object interface{}) error {
	return lruCache.Store.Find(key, object)
}

func (lruCache *LRUCache) Close() {
	lruCache.Store.Close()
}

func NewLRUCache(itemsLimit uint, store DataStore) (*LRUCache, error) {
	if itemsLimit == 0 {
		return nil, errors.New("Items limit must be greater 0")
	}

	return &LRUCache{
		ItemsLimit: itemsLimit,
		Store:      store,
	}, nil
}
