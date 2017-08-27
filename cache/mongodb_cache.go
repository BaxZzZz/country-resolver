package cache

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type mongoDBStore struct {
	session    *mgo.Session
	collection *mgo.Collection
}

func (store *mongoDBStore) GetItemsCount() (uint, error) {
	count, err := store.collection.Find(bson.M{}).Count()
	if err != nil {
		return 0, nil
	}

	return uint(count), nil
}

func (store *mongoDBStore) SetItem(item *Item) error {
	return store.collection.Insert(item)
}

func (store *mongoDBStore) FindItem(key string) (*Item, error) {
	item := &Item{}
	err := store.collection.Find(bson.M{"key": key}).One(item)
	if err != nil {
		return nil, err
	}

	err = store.collection.Update(bson.M{"key": key}, bson.M{"$set": bson.M{"timestamp": makeTimestamp()}})
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (store *mongoDBStore) RemoveOldItem() error {
	return store.collection.Remove(bson.M{"timestamp": bson.M{"$lt": makeTimestamp()}})
}

func (store *mongoDBStore) Close() {
	store.session.Close()
}

func NewMongoDbStore(url string, dbName string, collection string) (*mongoDBStore, error) {
	session, err := mgo.Dial(url)
	if err != nil {
		return nil, err
	}

	session.SetMode(mgo.Monotonic, true)

	store := &mongoDBStore{
		session:    session,
		collection: session.DB(dbName).C(collection),
	}

	index := mgo.Index{
		Key:        []string{"key", "timestamp"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}

	err = store.collection.EnsureIndex(index)
	if err != nil {
		return nil, err
	}

	return store, nil
}
