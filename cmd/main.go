package main

import (
	"log"
	"net"
	"os"

	"github.com/tn259/cc-redis/database"
	"github.com/tn259/cc-redis/resp"
)

func main() {
	// Open log file
	lf, err := os.OpenFile("cc-redis.log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer lf.Close()
	log.SetOutput(lf)
	log.Printf("Starting cc-redis")

	// Init the database
	_ = database.Database()

	// Listen for client connections on port 6379
	listener, err := net.Listen("tcp", ":6379")
	if err != nil {
		log.Fatal("Error: net.Listen():", err)
		return
	}
	defer listener.Close()

	// Accept client connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error: listener.Accept():", err)
			continue
		}

		// Handle client connection in a separate goroutine
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	b := make([]byte, 1024)
	for {
		_, err := conn.Read(b)
		if err != nil {
			log.Println("Error: conn.Read():", err)
			return
		}
		log.Println("Received:", string(b))

		// Parse the command
		parser := &resp.CommandParser{}
		cmd, err := parser.Parse(string(b))
		if err != nil {
			log.Println("Error: parser.Parse():", err)
			rErr := &resp.Error{Prefix: "ERR", Message: err.Error()}
			conn.Write([]byte(rErr.Serialize()))
			continue
		}

		// Execute the command
		res, err := cmd.Execute()
		if err != nil {
			log.Println("Error: cmd.Execute():", err)
			rErr := &resp.Error{Prefix: "ERR", Message: err.Error()}
			conn.Write([]byte(rErr.Serialize()))
			continue
		}

		// Serialize the command response
		_, err = conn.Write([]byte(res.Serialize()))
		if err != nil {
			log.Println("Error: conn.Write():", err)
			continue
		}
	}
}
