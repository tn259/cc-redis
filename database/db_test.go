package database

import (
	"testing"
	"time"
)

func TestDatabase_SetAndGet(t *testing.T) {
	db := Database()

	// Test Set and Get
	key := "mykey"
	value := "myvalue"
	db.Set(key, value, nil)

	result, ok := db.Get(key)
	if !ok {
		t.Errorf("Expected key %s to exist in the database", key)
	}

	if result != value {
		t.Errorf("Expected value %s, got %s", value, result)
	}
}

func TestDatabase_SetWithExpiration_Expired(t *testing.T) {
	db := Database()

	// Test Set with expiration
	key := "mykey"
	value := "myvalue"
	expiry := time.Now().Add(5 * time.Second)
	db.Set(key, value, &expiry)

	result, ok := db.Get(key)
	if !ok {
		t.Errorf("Expected key %s to exist in the database", key)
	}

	if result != value {
		t.Errorf("Expected value %s, got %s", value, result)
	}

	// Wait for the expiration to happen
	time.Sleep(6 * time.Second)

	_, ok = db.Get(key)
	if ok {
		t.Errorf("Expected key %s to be expired", key)
	}
}

func TestDatabase_SetWithExpiration_NotExpired(t *testing.T) {
	db := Database()

	// Test Set with
	key := "mykey"
	value := "myvalue"
	expiry := time.Now().Add(5 * time.Minute)
	db.Set(key, value, &expiry)

	// Wait for a bit - expiration should not happen
	time.Sleep(6 * time.Second)

	result, ok := db.Get(key)
	if !ok {
		t.Errorf("Expected key %s to exist in the database", key)
	}

	if result != value {
		t.Errorf("Expected value %s, got %s", value, result)
	}
}

func TestDatabase_Delete(t *testing.T) {
	db := Database()

	// Test Delete
	key := "mykey"
	value := "myvalue"
	db.Set(key, value, nil)

	db.Delete(key)

	_, ok := db.Get(key)
	if ok {
		t.Errorf("Expected key %s to be deleted from the database", key)
	}
}
