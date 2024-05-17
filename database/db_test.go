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

	keys := []string{key, "nonexisitentkey"}
	c := db.Delete(keys)
	if c != 1 {
		t.Errorf("Expected 1 key to be deleted, got %d", c)
	}

	_, ok := db.Get(key)
	if ok {
		t.Errorf("Expected key %s to be deleted from the database", key)
	}
}

func TestDatabase_ListLPush(t *testing.T) {
	db := Database()

	// Test ListLPush
	key := "mylist"
	value1 := "myvalue"
	db.ListLPush(key, value1)
	value2 := "myvalue2"
	db.ListLPush(key, value2)

	result, err := db.ListRange(key, "0", "1")
	if err != nil {
		t.Errorf("Expected key %s to exist in the database", key)
	}

	if len(result) != 2 {
		t.Errorf("Expected list length to be 1, got %d", len(result))
	}

	if result[0] != value2 {
		t.Errorf("Expected value %s, got %s", value2, result[0])
	}
	if result[1] != value1 {
		t.Errorf("Expected value %s, got %s", value1, result[1])
	}
}

func TestDatabase_ListRPush(t *testing.T) {
	db := Database()

	// Test ListRPush
	key := "mylist"
	value1 := "myvalue"
	db.ListRPush(key, value1)
	value2 := "myvalue2"
	db.ListRPush(key, value2)

	result, err := db.ListRange(key, "0", "1")
	if err != nil {
		t.Errorf("Expected key %s to exist in the database", key)
	}

	if len(result) != 2 {
		t.Errorf("Expected list length to be 1, got %d", len(result))
	}

	if result[0] != value1 {
		t.Errorf("Expected value %s, got %s", value1, result[0])
	}
	if result[1] != value2 {
		t.Errorf("Expected value %s, got %s", value2, result[1])
	}
}

func TestDatabase_ListRange(t *testing.T) {
	db := Database()

	// Test ListRange
	key := "mylist"
	value1 := "myvalue"
	db.ListRPush(key, value1)
	value2 := "myvalue2"
	db.ListRPush(key, value2)
	value3 := "myvalue3"
	db.ListRPush(key, value3)

	type startStop struct {
		start, stop string
	}

	for ss, expected := range map[startStop][]string{
		{"0", "0"}:   {value1},
		{"0", "1"}:   {value1, value2},
		{"0", "-1"}:  {value1, value2, value3},
		{"-1", "-1"}: {value3},
		{"-2", "-1"}: {value2, value3},
		{"-3", "-1"}: {value1, value2, value3},
		{"-3", "-2"}: {value1, value2},
		{"-3", "-3"}: {value1},
		{"-4", "-3"}: {value1},
		{"4", "2"}:   {},
		{"1", "-4"}:  {},
	} {
		result, err := db.ListRange(key, ss.start, ss.stop)
		if err != nil {
			t.Errorf("Expected key %s to exist in the database", key)
		}

		if len(result) != len(expected) {
			t.Errorf("Expected list length to be %d, got %d", len(expected), len(result))
		}

		for i := range result {
			if result[i] != expected[i] {
				t.Errorf("Expected value %s, got %s", expected[i], result[i])
			}
		}
	}
}
