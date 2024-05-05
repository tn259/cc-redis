package resp

import (
	"fmt"

	"github.com/tn259/cc-redis/database"
)

// https://redis.io/docs/latest/develop/reference/protocol-spec/#sending-commands-to-a-redis-server

// Command represents a Redis command
type Command interface {
	// Execute the command
	// Returns the command response and an error
	Execute() (Type, error)
}

type Parser interface {
	Parse(string) (Command, error)
}

// CommandParser is a parser for Redis commands
type CommandParser struct {
}

func (*CommandParser) Parse(input string) (Command, error) {
	// Commands are RESP arrays of RESP bulk strings
	// Parse the input as an RESP array
	a := &Array{}
	if err := a.Deserialize(input); err != nil {
		return nil, fmt.Errorf("failed to parse command: %v", err)
	}

	// The first element of the array is the command name
	if len(a.Elements) == 0 {
		return nil, fmt.Errorf("missing command name")
	}

	arg0 := a.Elements[0].(*BulkString)
	var arg1 *BulkString
	if len(a.Elements) > 1 {
		arg1 = a.Elements[1].(*BulkString)
	}

	// Create a new command based on the command name
	switch cmd := arg0.Value; cmd {
	case "PING":
		return &Ping{arg: arg1}, nil
	case "ECHO":
		return &Echo{arg: arg1}, nil
	case "SET":
		return NewSet(a)
	case "GET":
		return &Get{key: arg1}, nil
	default:
		return nil, fmt.Errorf("unknown command: %s", cmd)
	}
}

// https://redis.io/docs/latest/commands/ping/
type Ping struct {
	arg *BulkString
}

func (p *Ping) Execute() (Type, error) {
	if p.arg == nil {
		return &SimpleString{Value: "PONG"}, nil
	}
	return p.arg, nil
}

// https://redis.io/docs/latest/commands/echo/
type Echo struct {
	arg *BulkString
}

func (e *Echo) Execute() (Type, error) {
	return e.arg, nil
}

// https://redis.io/docs/latest/commands/set/
type Set struct {
	key   *BulkString
	value *BulkString
}

func NewSet(a *Array) (*Set, error) {
	if len(a.Elements) != 3 {
		return nil, fmt.Errorf("SET command requires 3 arguments")
	}
	key := a.Elements[1].(*BulkString)
	value := a.Elements[2].(*BulkString)
	return &Set{key: key, value: value}, nil
}

func (s *Set) Execute() (Type, error) {
	database.StringsDB().Set(s.key.Value, s.value.Value)
	return &SimpleString{Value: "OK"}, nil
}

// https://redis.io/docs/latest/commands/get/
type Get struct {
	key *BulkString
}

func (g *Get) Execute() (Type, error) {
	value, ok := database.StringsDB().Get(g.key.Value)
	if !ok {
		return &BulkString{IsNull: true}, nil
	}
	return &BulkString{Value: value}, nil
}
