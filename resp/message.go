package resp

import (
	"fmt"
	"strconv"
	"strings"
)

const CRLF = "\r\n"

// Type is an interface for RESP types
type Type interface {
	Serialize() string
	Deserialize(string) error
}

// SimpleString is a RESP simple string
type SimpleString struct {
	Value string
}

// Implement the Type interface for SimpleString
func (s *SimpleString) Serialize() string {
	return "+" + s.Value + CRLF
}

func (s *SimpleString) Deserialize(input string) error {
	if input[0] != '+' {
		return fmt.Errorf("invalid RESP simple string")
	}

	firstCRLF := strings.Index(input, CRLF)
	if firstCRLF == -1 {
		return fmt.Errorf("invalid RESP simple string - no CRLF found after value")
	}

	s.Value = input[1:firstCRLF]
	return nil
}

// Error is a RESP error
type Error struct {
	Prefix  string
	Message string
}

func (e *Error) Serialize() string {
	return "-" + e.Prefix + " " + e.Message + CRLF
}

func (e *Error) Deserialize(input string) error {
	if input[0] != '-' {
		return fmt.Errorf("invalid RESP error")
	}

	firstSpace := strings.Index(input, " ")
	if firstSpace == -1 {
		return fmt.Errorf("invalid RESP error - no space found between prefix and message")
	}
	e.Prefix = input[1:firstSpace]

	firstCRLF := strings.Index(input, CRLF)
	if firstCRLF == -1 {
		return fmt.Errorf("invalid RESP error - no CRLF found after message")
	}
	e.Message = input[firstSpace+1 : firstCRLF]
	return nil
}

// Integer is a RESP integer
type Integer struct {
	Value int
}

// Implement the Type interface for Integer
func (i *Integer) Serialize() string {
	return ":" + strconv.Itoa(i.Value) + CRLF
}

func (i *Integer) Deserialize(input string) error {
	if input[0] != ':' {
		return fmt.Errorf("invalid RESP integer")
	}

	firstCRLF := strings.Index(input, CRLF)
	if firstCRLF == -1 {
		return fmt.Errorf("invalid RESP integer - no CRLF found after value")
	}

	value, err := strconv.Atoi(input[1:firstCRLF])
	if err != nil {
		return fmt.Errorf("invalid RESP integer strconv.Atoi(): %v", err)
	}

	i.Value = value
	return nil
}

// And so on for BulkString, and Array...

// BulkString is a RESP bulk string
// Value can be binary data
type BulkString struct {
	Value  string
	IsNull bool
}

func (b *BulkString) Serialize() string {
	if b.IsNull {
		return "$-1" + CRLF
	}
	return "$" + strconv.Itoa(len(b.Value)) + CRLF + b.Value + CRLF
}

func (b *BulkString) Deserialize(input string) error {
	if input[0] != '$' {
		return fmt.Errorf("invalid RESP bulk string")
	}

	firstCRLF := strings.Index(input, CRLF)
	if firstCRLF == -1 {
		return fmt.Errorf("invalid RESP bulk string - no CRLF found after length")
	}

	// null case
	if input[1:firstCRLF] == "-1" {
		b.Value = ""
		b.IsNull = true
		return nil
	}

	length, err := strconv.Atoi(input[1:firstCRLF])
	if err != nil {
		return fmt.Errorf("invalid RESP bulk string length - %v", err)
	}

	b.Value = input[firstCRLF+2 : firstCRLF+2+length]
	return nil
}

// Array is a RESP array
type Array struct {
	Elements []Type
	IsNull   bool
}

func (a *Array) Serialize() string {
	if a.IsNull {
		return "*-1" + CRLF
	}
	var result string
	result += "*" + strconv.Itoa(len(a.Elements)) + CRLF
	for _, element := range a.Elements {
		result += element.Serialize()
	}
	return result
}

func (a *Array) Deserialize(input string) error {
	if input[0] != '*' {
		return fmt.Errorf("invalid RESP array")
	}

	firstCRLF := strings.Index(input, CRLF)
	if firstCRLF == -1 {
		return fmt.Errorf("invalid RESP array - no CRLF found after length")
	}

	// null case
	if input[1:firstCRLF] == "-1" {
		a.Elements = nil
		a.IsNull = true
		return nil
	}

	length, err := strconv.Atoi(input[1:firstCRLF])
	if err != nil {
		return fmt.Errorf("invalid RESP array length - %v", err)
	}

	remaining := input[firstCRLF+2:]
	a.Elements = make([]Type, length)
	for i := 0; i < length; i++ {
		var element Type
		switch remaining[0] {
		case '+':
			element = &SimpleString{}
		case '-':
			element = &Error{}
		case ':':
			element = &Integer{}
		case '$':
			element = &BulkString{}
		case '*':
			element = &Array{}
		default:
			return fmt.Errorf("invalid RESP type")
		}

		err = element.Deserialize(remaining)
		if err != nil {
			return err
		}

		a.Elements[i] = element
		remaining = remaining[len(element.Serialize()):]
	}
	return nil
}
