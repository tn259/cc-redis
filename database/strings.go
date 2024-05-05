package database

import "sync"

// stringsDB represents a Redis in-memory strings database.
type stringsDB struct {
	data map[string]string
	mu   sync.RWMutex
}

var db *stringsDB
var once sync.Once

// NewStringsDB creates a new instance of StringsDB.
func StringsDB() *stringsDB {
	once.Do(func() {
		db = &stringsDB{
			data: make(map[string]string),
		}
	})
	return db
}

// Set sets the value of a key in the database.
func (db *stringsDB) Set(key, value string) {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.data[key] = value
}

// Get retrieves the value of a key from the database.
func (db *stringsDB) Get(key string) (string, bool) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	value, ok := db.data[key]
	return value, ok
}

// Delete deletes a key from the database.
func (db *stringsDB) Delete(key string) {
	db.mu.Lock()
	defer db.mu.Unlock()
	delete(db.data, key)
}
