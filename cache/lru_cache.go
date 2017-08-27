package cache

import "github.com/mitchellh/mapstructure"

type LruCache struct {
	ItemsLimit uint
	Store      DataStore
}

func (lruCache *LruCache) Set(key string, object interface{}) error {
	countItems, err := lruCache.Store.GetItemsCount()
	if err != nil {
		return err
	}

	if countItems >= lruCache.ItemsLimit {
		err = lruCache.Store.RemoveOldItem()
		if err != nil {
			return err
		}
	}

	lruCache.Store.SetItem(&Item{
		Key:       key,
		Object:    object,
		Timestamp: makeTimestamp(),
	})

	return nil
}

func (lruCache *LruCache) Get(key string, object interface{}) error {
	item, err := lruCache.Store.FindItem(key)
	if err != nil {
		return err
	}

	return mapstructure.Decode(item.Object, object)
}

func (lruCache *LruCache) Close() {
	lruCache.Store.Close()
}
