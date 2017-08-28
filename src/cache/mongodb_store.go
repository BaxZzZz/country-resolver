package cache

import (
	"time"

	"github.com/mitchellh/mapstructure"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
)

// Item with data and update timestamp
type Item struct {
	Key       string
	Object    interface{}
	Timestamp int64
}

// Data store on MongoDB database
type MongoDBStore struct {
	session    *mgo.Session
	collection *mgo.Collection
}

// Make new timestamp
func makeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

// Get objects amount
func (store *MongoDBStore) GetCount() (uint, error) {
	count, err := store.collection.Find(bson.M{}).Count()
	if err != nil {
		return 0, nil
	}

	return uint(count), nil
}

// Set new object to store by key
func (store *MongoDBStore) Set(key string, object interface{}) error {
	item := &Item{
		Key:       key,
		Object:    object,
		Timestamp: makeTimestamp(),
	}

	return store.collection.Insert(item)
}

// Find object by key value
func (store *MongoDBStore) Find(key string, object interface{}) error {
	item := &Item{
		Key:       key,
		Object:    object,
		Timestamp: makeTimestamp(),
	}

	err := store.collection.Find(bson.M{"key": key}).One(item)
	if err != nil {
		return err
	}

	err = mapstructure.Decode(item.Object, object)
	if err != nil {
		return err
	}

	err = store.collection.Update(bson.M{"key": key}, bson.M{"$set": bson.M{"timestamp": makeTimestamp()}})
	if err != nil {
		return err
	}

	return nil
}

// Remove the oldest record
func (store *MongoDBStore) RemoveOld() error {
	return store.collection.Remove(bson.M{"timestamp": bson.M{"$lt": makeTimestamp()}})
}

// Close working with store
func (store *MongoDBStore) Close() {
	store.session.Close()
}

// Creates new MongoDB store instance
func NewMongoDBStore(url string, dbName string, username string, password string,
	collection string) (*MongoDBStore, error) {

	log.Println("Connect to DB, url: " + url + ", dbName: " + dbName + ", collection: " + collection)

	dialInfo := &mgo.DialInfo{
		Addrs:    []string{url},
		Timeout:  10 * time.Second,
		Database: dbName,
		Username: username,
		Password: password,
	}

	session, err := mgo.DialWithInfo(dialInfo)

	if err != nil {
		return nil, err
	}

	session.SetMode(mgo.Monotonic, true)

	store := &MongoDBStore{
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
