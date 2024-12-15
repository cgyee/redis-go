package main

import (
	"sync"
	"time"
)

var once sync.Once
var instance *Store

type Data struct {
	Key    string
	Type   string
	Value  []byte
	Expire time.Time
}

type Store struct {
	data map[string]Data
	mu   sync.Mutex
}

func (store *Store) Get(key string) (row Data) {
	store.mu.Lock()
	row, ok := store.data[key]
	defer store.mu.Unlock()

	if !ok {
		return Data{}
	}
	if time.Now().After(row.Expire) {
		delete(store.data, key)
		return Data{}
	}
	return row
}

func (store *Store) Set(key string, t string, v []byte) (ok bool) {
	store.mu.Lock()
	_, ok = store.data[key]
	defer store.mu.Unlock()

	if ok {
		return false
	}
	store.data[key] = Data{Key: key, Type: t, Value: v}
	return true
}

func (store *Store) SetEx(key string, t string, v []byte, tm int64) (ok bool) {
	store.mu.Lock()
	_, ok = store.data[key]
	defer store.mu.Unlock()

	if ok {
		return false
	}
	duration := time.Millisecond * time.Duration(tm)
	expire := time.Now().Add(duration)
	store.data[key] = Data{Key: key, Type: t, Value: v, Expire: expire}
	return true
}

func GetInstance() *Store {
	once.Do(func() {
		instance = &Store{data: map[string]Data{}}
	})
	return instance
}
