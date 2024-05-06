package database

import (
	"sync"
	"time"
)

type dbstring struct {
	value  string
	expiry *time.Time
}

// stringsDB represents a Redis in-memory strings database.
type stringsDB struct {
	data map[string]dbstring
	mu   sync.RWMutex
}

var db *stringsDB
var once sync.Once

// NewStringsDB creates a new instance of StringsDB.
func StringsDB() *stringsDB {
	once.Do(func() {
		db = &stringsDB{
			data: make(map[string]dbstring),
		}
	})
	return db
}

// Set sets the value of a key in the database.
func (db *stringsDB) Set(key, value string, expiry *time.Time) {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.data[key] = dbstring{value: value, expiry: expiry}
}

// Get retrieves the value of a key from the database.
func (db *stringsDB) Get(key string) (string, bool) {
	db.mu.RLock()
	s, ok := db.data[key]
	db.mu.RUnlock()
	if ok {
		if s.expiry != nil && s.expiry.Before(time.Now()) {
			// Passive expiry
			// TODO implement active expiry - https://redis.io/commands/expire
			db.Delete(key)
			return "", false
		}
	}
	return s.value, ok
}

// Delete deletes a key from the database.
func (db *stringsDB) Delete(key string) {
	db.mu.Lock()
	defer db.mu.Unlock()
	delete(db.data, key)
}
