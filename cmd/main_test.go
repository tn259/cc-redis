package main

import (
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/go-redis/redis"
)

var cmd *exec.Cmd

func start(t *testing.T, client *redis.Client, startup bool) {
	// Start your Redis server...
	if startup {
		cmd = exec.Command("go", "run", "main.go")
		if err := cmd.Start(); err != nil {
			t.Fatalf("Could not start Redis server: %v", err)
		}
	}
	for {
		_, err := client.Ping().Result()
		if err == nil {
			t.Log("Redis server is ready")
			break
		}
		t.Log("Waiting for Redis server to start...")
		time.Sleep(2 * time.Second)
	}
}
func stop(t *testing.T) {
	// Stop your Redis server...
	if cmd == nil || cmd.Process == nil {
		t.Fatalf("Could not stop Redis server: process is nil")
	}
	if err := cmd.Process.Signal(os.Interrupt); err != nil {
		t.Fatalf("Could not stop Redis server: %v", err)
		err := cmd.Process.Kill()
		if err != nil {
			t.Fatalf("Could not kill Redis server: %v", err)
		}
	}
	// Wait for the process to exit
	cmd.Wait()
}

func SetTest(t *testing.T, client *redis.Client) {
	// Set a key-value pair
	err := client.Set("key", "value", 0).Err()
	if err != nil {
		t.Fatalf("Could not set key-value pair: %v", err)
	}
}

func NilGetTest(t *testing.T, client *redis.Client) {
	// Get a key that does not exist
	value, err := client.Get("unknownkey").Result()
	if err != redis.Nil {
		t.Fatalf("Expected key to be nil: %v", err)
	}
	if value != "" {
		t.Fatalf("Expected value to be empty: %v", value)
	}
}

func GetTest(t *testing.T, client *redis.Client) {
	// Set a key-value pair
	err := client.Set("key", "value", 0).Err()
	if err != nil {
		t.Fatalf("Could not set key-value pair: %v", err)
	}

	// Get a key that exists
	value, err := client.Get("key").Result()
	if err != nil {
		t.Fatalf("Could not get key-value pair: %v", err)
	}
	if value != "value" {
		t.Fatalf("Expected value to be 'value': %v", value)
	}
}

func ExistsTest(t *testing.T, client *redis.Client) {
	// 'key2' does not exist in the db so now set it
	err := client.Set("key2", "value1", 0).Err()
	if err != nil {
		t.Fatalf("Could not set key-value pair: %v", err)
	}

	// 'key' already exists in the db
	keys := []string{"key10", "key2"}
	cmd := client.Exists(keys...)
	if cmd.Err() != nil {
		t.Fatalf("Could not check if key exists: %v", cmd.Err())
	}
	if cmd.Val() != 1 {
		t.Fatalf("Expected key to exist: %v", cmd.Val())
	}
}

func DeleteTest(t *testing.T, client *redis.Client) {
	// Set a key-value pair
	err := client.Set("key", "value", 0).Err()
	if err != nil {
		t.Fatalf("Could not set key-value pair: %v", err)
	}

	// Check key exists
	cmd := client.Exists("key")
	if cmd.Err() != nil {
		t.Fatalf("Could not check if key exists: %v", cmd.Err())
	}
	if cmd.Val() != 1 {
		t.Fatalf("Expected key to exist: %v", cmd.Val())
	}

	// Delete a key that exists
	cmd = client.Del([]string{"key", "key10"}...)
	if cmd.Err() != nil {
		t.Fatalf("Could not delete key: %v", cmd.Err())
	}

	// Check key does not exist
	if cmd.Val() != 1 {
		t.Fatalf("Expected key to be deleted: %v", cmd.Val())
	}
}

func IncrTest(t *testing.T, client *redis.Client) {
	// Set a key-value pair
	err := client.Set("stringkey", "qwerty", 0).Err()
	if err != nil {
		t.Fatalf("Could not set key-value pair: %v", err)
	}
	err = client.Incr("stringkey").Err()
	if err == nil {
		t.Fatalf("Expected error: %v", err)
	}

	// Increment a key that does not exist
	cmd := client.Incr("nonexistentkey")
	if cmd.Err() != nil {
		t.Fatalf("Could not increment key: %v", cmd.Err())
	}
	if cmd.Val() != 1 {
		t.Fatalf("Expected value to be 1: %v", cmd.Val())
	}

	// Increment a key that exists
	err = client.Set("intkey", "23", 0).Err()
	if err != nil {
		t.Fatalf("Could not set key-value pair: %v", err)
	}
	cmd = client.Incr("intkey")
	if cmd.Err() != nil {
		t.Fatalf("Could not increment key: %v", cmd.Err())
	}
	if cmd.Val() != 24 {
		t.Fatalf("Expected value to be24: %v", cmd.Val())
	}
}

func DecrTest(t *testing.T, client *redis.Client) {
	// Set a key-value pair
	err := client.Set("stringkey2", "qwerty", 0).Err()
	if err != nil {
		t.Fatalf("Could not set key-value pair: %v", err)
	}
	err = client.Decr("stringkey2").Err()
	if err == nil {
		t.Fatalf("Expected error: %v", err)
	}

	// Decrement a key that does not exist
	cmd := client.Decr("nonexistentkey2")
	if cmd.Err() != nil {
		t.Fatalf("Could not decrement key: %v", cmd.Err())
	}
	if cmd.Val() != -1 {
		t.Fatalf("Expected value to be -1: %v", cmd.Val())
	}

	// Decrement a key that exists
	err = client.Set("intkey2", "23", 0).Err()
	if err != nil {
		t.Fatalf("Could not set key-value pair: %v", err)
	}
	cmd = client.Decr("intkey2")
	if cmd.Err() != nil {
		t.Fatalf("Could not decrement key: %v", cmd.Err())
	}
	if cmd.Val() != 22 {
		t.Fatalf("Expected value to be 22: %v", cmd.Val())
	}
}

func TestRedisCommands(t *testing.T) {
	// Define the commands to be sent during the test
	tests := []struct {
		name string
		test func(t *testing.T, client *redis.Client)
	}{
		{name: "Set", test: SetTest},
		{name: "Nil Get", test: NilGetTest},
		{name: "Get", test: GetTest},
		{name: "Exists", test: ExistsTest},
		{name: "Delete", test: DeleteTest},
		{name: "Incr", test: IncrTest},
		{name: "Decr", test: DecrTest},
		// Add more commands here...
	}

	// Create a Redis client
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // address of your Redis server
		Password: "",               // no password
		DB:       0,                // use default DB
	})
	// Start the Redis server
	startRedis := false
	start(t, client, startRedis)
	if startRedis {
		defer stop(t)
	}
	// Execute the test function for each command
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Execute the test function
			tt.test(t, client)
		})
	}
	// Close the client
	client.Close()
}
