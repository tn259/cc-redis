package main

import (
	"bytes"
	"os"
	"os/exec"
	"syscall"
	"testing"
	"time"

	"github.com/go-redis/redis"
	"github.com/tn259/cc-redis/database"
)

var cmd *exec.Cmd

func newClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}
func start(t *testing.T, client *redis.Client, startup bool) {
	// Start your Redis server...
	if startup {
		cmd = exec.Command("go", "run", "main.go")
		// Create a new process group to terminate child processes
		cmd.SysProcAttr = &syscall.SysProcAttr{
			Setpgid: true,
		}
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
	pgid, err := syscall.Getpgid(cmd.Process.Pid)
	if err != nil {
		t.Fatalf("Could not get pgid of Redis Server: %v", err)
	}
	if err := syscall.Kill(-pgid, syscall.SIGKILL); err != nil {
		t.Fatalf("Could not stop Redis server: %v", err)
	}
}
func fileContentsEqual(t *testing.T, file1, file2 string) bool {
	// The files we are testing are not intended to be large
	b1, err := os.ReadFile(file1)
	if err != nil {
		t.Fatalf("Could not read file: %v", err)
	}
	b2, err := os.ReadFile(file2)
	if err != nil {
		t.Fatalf("Could not read file: %v", err)
	}

	return bytes.Equal(b1, b2)
}
func copyFile(t *testing.T, src, dst string) {
	b, err := os.ReadFile(src)
	if err != nil {
		t.Fatalf("Could not read file: %v", err)
	}
	if err := os.WriteFile(dst, b, 0666); err != nil {
		t.Fatalf("Could not write file: %v", err)
	}
}
func removeFile(t *testing.T, file string) {
	if err := os.Remove(file); err != nil {
		t.Fatalf("Could not remove file: %v", err)
	}
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

func ListTest(t *testing.T, client *redis.Client) {
	// Set a key-value pair
	err := client.LPush("mylist", "value1").Err()
	if err != nil {
		t.Fatalf("Could not set key-value pair: %v", err)
	}
	err = client.LPush("mylist", "value2").Err()
	if err != nil {
		t.Fatalf("Could not set key-value pair: %v", err)
	}
	err = client.RPush("mylist", "value3").Err()
	if err != nil {
		t.Fatalf("Could not set key-value pair: %v", err)
	}

	// Get a range of values
	cmd := client.LRange("mylist", 0, -1)
	if cmd.Err() != nil {
		t.Fatalf("Could not get range of values: %v", cmd.Err())
	}
	if len(cmd.Val()) != 3 {
		t.Fatalf("Expected list length to be 3: %v", len(cmd.Val()))
	}
	if cmd.Val()[0] != "value2" {
		t.Fatalf("Expected value to be 'value2': %v", cmd.Val()[0])
	}
	if cmd.Val()[1] != "value1" {
		t.Fatalf("Expected value to be 'value1': %v", cmd.Val()[1])
	}
	if cmd.Val()[2] != "value3" {
		t.Fatalf("Expected value to be 'value3': %v", cmd.Val()[2])
	}
}

func SaveTest(t *testing.T, client *redis.Client) {
	// At this point the redis should have the following data...
	// Save the database
	cmd := client.Save()
	if cmd.Err() != nil {
		t.Fatalf("Could not save the database: %v", cmd.Err())
	}
	defer removeFile(t, database.RDBFilename)
	// Direct file comparison
	if !fileContentsEqual(t, database.RDBFilename, database.RDBFilename+".test") {
		t.Fatalf("Expected files to be equal")
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
		{name: "ListTest", test: ListTest},
		{name: "Save", test: SaveTest},
		// Add more commands here...
	}

	// Create a Redis client
	client := newClient()
	// Start the Redis server
	startRedis := true
	start(t, client, startRedis)
	if startRedis {
		defer func() {
			client.Close()
			stop(t)
		}()
	}
	// Execute the test function for each command
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Execute the test function
			tt.test(t, client)
		})
	}
}

func TestRedisCommands_ReadOnStartup(t *testing.T) {
	// Copy the test RDB file to the actual RDB file
	copyFile(t, database.RDBFilename+".test", database.RDBFilename)
	// Create client
	client := newClient()
	// Start the Redis server
	startRedis := true
	start(t, client, startRedis)
	if startRedis {
		defer func() {
			client.Close()
			stop(t)
			removeFile(t, database.RDBFilename)
		}()
	}
	// Check if the data is loaded
	mylistKey := client.LRange("mylist", 0, -1)
	if mylistKey.Err() != nil {
		t.Fatalf("Could not get range of values: %v", mylistKey.Err())
	}
	if len(mylistKey.Val()) != 3 {
		t.Fatalf("Expected list length to be 3: %v", len(mylistKey.Val()))
	}

	// TODO check the rest...
}
