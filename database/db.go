package database

import (
	"fmt"
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
	db.mu.RUnlock()
	if !ok {
		return "", false
	}
	s, ok := e.(dbstring)
	if !ok {
		return "", false
	}
	if s.expiry != nil && s.expiry.Before(time.Now()) {
		// Passive expiry
		// TODO implement active expiry - https://redis.io/commands/expire
		db.Delete(key)
		return "", false
	}
	return s.value, ok
}

// Delete deletes a key from the database.
func (db *DB) Delete(key string) {
	db.mu.Lock()
	defer db.mu.Unlock()
	delete(db.data, key)
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
