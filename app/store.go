package main

import (
	"fmt"
	"sync"
)

var once sync.Once
var instance *Store

type Data struct {
	Key   string
	Type  string
	Value []byte
}

type Store struct {
	data map[string]Data
}

func (store *Store) Get(key string) (row Data) {
	row, ok := store.data[key]

	if !ok {
		return Data{}
	}
	return row
}

func (store *Store) Set(key string, t string, v []byte) (ok bool) {

	_, ok = store.data[key]
	if ok {
		return false
	}
	fmt.Println("Set v: ", string(v))
	store.data[key] = Data{Key: key, Type: t, Value: v}
	return true
}
func GetInstance() *Store {
	once.Do(func() {
		instance = &Store{data: map[string]Data{}}
	})
	return instance
}
