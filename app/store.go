package main

import (
	"fmt"
	"sync"
	"time"
)

var once sync.Once
var instance *Store

type Data struct {
	Key     string
	Type    string
	Value   []byte
	Expires bool
	Expire  time.Time
}

type Store struct {
	data map[string]Data
	mu   sync.Mutex
}

func (store *Store) Get(key string) (row Data) {
	store.mu.Lock()
	row, ok := store.data[key]
	fmt.Println(row.Key, string(row.Value), ok)
	defer store.mu.Unlock()

	if !ok {
		return Data{}
	}
	if row.Expires && time.Now().After(row.Expire) {
		delete(store.data, key)
		return Data{}
	}
	return row
}

func (store *Store) Set(key string, t string, v []byte) (ok bool) {
	store.mu.Lock()
	defer store.mu.Unlock()
	store.data[key] = Data{Key: key, Type: t, Value: v, Expires: false}
	return true
}

func (store *Store) SetEx(key string, t string, v []byte, tm int64) (ok bool) {
	store.mu.Lock()
	_, ok = store.data[key]
	defer store.mu.Unlock()
	duration := time.Millisecond * time.Duration(tm)
	expire := time.Now().Add(duration)
	store.data[key] = Data{Key: key, Type: t, Value: v, Expire: expire, Expires: true}
	return true
}

func GetInstance() *Store {
	once.Do(func() {
		instance = &Store{data: map[string]Data{}}
	})
	return instance
}
