package main

import (
	"log"
	"net"
	"os"

	"github.com/tn259/cc-redis/database"
	"github.com/tn259/cc-redis/resp"
)

type Command struct {
	cmd  resp.Command
	conn *net.Conn
}

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

	commandChan := make(chan *Command)

	// Accept client connections
	go func(c chan *Command) {
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Println("Error: listener.Accept():", err)
				continue
			}

			// Handle client connection in a separate goroutine
			go handleConnection(conn, c)
		}
	}(commandChan)

	for c := range commandChan {
		handleCommand(c)
	}
}

func handleConnection(conn net.Conn, commandChan chan *Command) {
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

		// Send the command to the command channel
		commandChan <- &Command{cmd: cmd, conn: &conn}
	}
}

func handleCommand(c *Command) {
	// Execute the command
	res, err := c.cmd.Execute()
	if err != nil {
		log.Println("Error: cmd.Execute():", err)
		rErr := &resp.Error{Prefix: "ERR", Message: err.Error()}
		_, err := (*c.conn).Write([]byte(rErr.Serialize()))
		if err != nil {
			log.Println("Error: conn.Write():", err)
		}
		return
	}

	// Serialize the command response
	_, err = (*c.conn).Write([]byte(res.Serialize()))
	if err != nil {
		log.Println("Error: conn.Write():", err)
	}
}
