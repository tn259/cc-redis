package database

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

type dbstring struct {
	value  string
	expiry *time.Time
}

type node struct {
	value string
	prev  *node
	next  *node
}

type dblist struct {
	head *node
	tail *node
}

// db represents a Redis in-memory strings database.
type DB struct {
	data map[string]interface{}
	mu   sync.RWMutex
}

var db *DB
var once sync.Once

// NewStringsDB creates a new instance of StringsDB.
func Database() *DB {
	once.Do(func() {
		db = &DB{
			data: make(map[string]interface{}),
		}
	})
	return db
}

// Set sets the value of a key in the database.
func (db *DB) Set(key, value string, expiry *time.Time) {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.data[key] = dbstring{value: value, expiry: expiry}
}

// Get retrieves the value of a key from the database.
func (db *DB) Get(key string) (string, bool) {
	db.mu.RLock()
	e, ok := db.data[key]
	if !ok {
		db.mu.RUnlock()
		return "", false
	}
	s, ok := e.(dbstring)
	if !ok {
		db.mu.RUnlock()
		return "", false
	}
	if s.expiry != nil && s.expiry.Before(time.Now()) {
		db.mu.RUnlock()
		// Passive expiry
		// TODO implement active expiry - https://redis.io/commands/expire
		db.Delete([]string{key})
		return "", false
	}
	db.mu.RUnlock()
	return s.value, ok
}

// Delete deletes a key from the database.
func (db *DB) Delete(keys []string) int {
	db.mu.Lock()
	defer db.mu.Unlock()
	c := 0
	for _, key := range keys {
		_, ok := db.data[key]
		if !ok {
			continue
		}
		delete(db.data, key)
		c++
	}
	return c
}

// ListLPush adds an element to the head of a list.
func (db *DB) ListLPush(key, value string) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	e, ok := db.data[key]
	if !ok {
		l := &dblist{}
		n := &node{value: value}
		l.head = n
		l.tail = n
		db.data[key] = l
		return nil
	}
	l, ok := e.(*dblist)
	if !ok {
		return fmt.Errorf("key %s does not contain a list", key)
	}
	n := &node{value: value}
	n.next = l.head
	l.head.prev = n
	l.head = n
	return nil
}

// ListRPush adds an element to the tail of a list.
func (db *DB) ListRPush(key, value string) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	e, ok := db.data[key]
	if !ok {
		l := &dblist{}
		n := &node{value: value}
		l.head = n
		l.tail = n
		db.data[key] = l
		return nil
	}
	l, ok := e.(*dblist)
	if !ok {
		return fmt.Errorf("key %s does not contain a list", key)
	}
	n := &node{value: value}
	n.prev = l.tail
	l.tail.next = n
	l.tail = n
	return nil
}

func (db *DB) ListRange(key, start, stop string) ([]string, error) {
	startInt, err := strconv.Atoi(start)
	if err != nil {
		return nil, fmt.Errorf("invalid start index %s", start)
	}
	stopInt, err := strconv.Atoi(stop)
	if err != nil {
		return nil, fmt.Errorf("invalid stop index %s", stop)
	}
	db.mu.RLock()
	defer db.mu.RUnlock()
	e, ok := db.data[key]
	if !ok {
		return nil, fmt.Errorf("key %s does not exist", key)
	}
	l, ok := e.(*dblist)
	if !ok {
		return nil, fmt.Errorf("key %s does not contain a list", key)
	}
	values := []string{}
	listLen := 0
	for n := l.head; n != nil; n = n.next {
		listLen++
	}
	if startInt < 0 {
		startInt = listLen + startInt
	}
	if startInt < 0 {
		startInt = 0
	}
	if startInt >= listLen {
		return values, nil
	}
	if stopInt < 0 {
		stopInt = listLen + stopInt
	}
	if stopInt < 0 {
		return values, nil
	}
	if stopInt >= listLen {
		stopInt = listLen - 1
	}
	for i, n := 0, l.head; n != nil; i, n = i+1, n.next {
		if i >= startInt && i <= stopInt {
			values = append(values, n.value)
		}
	}
	return values, nil
}

// Write rdb file
func (db *DB) Save() error {
	writer := &RDBWriter{db: db}
	return writer.Write()
}
