package resp

import (
	"testing"
)

func TestSimpleString_Serialize(t *testing.T) {
	s := &SimpleString{Value: "Hello World"}
	expected := "+Hello World\r\n"
	if got := s.Serialize(); got != expected {
		t.Errorf("Serialize() = %s; want %s", got, expected)
	}
}

func TestSimpleString_Deserialize(t *testing.T) {
	s := &SimpleString{}
	input := "+Hello World\r\n"
	if err := s.Deserialize(input); err != nil {
		t.Errorf("Deserialize() returned an error: %v", err)
	}
	expected := "Hello World"
	if s.Value != expected {
		t.Errorf("Deserialize() = %s; want %s", s.Value, expected)
	}
}

func TestError_Serialize(t *testing.T) {
	e := &Error{Prefix: "ERR", Message: "Something went wrong"}
	expected := "-ERR Something went wrong\r\n"
	if got := e.Serialize(); got != expected {
		t.Errorf("Serialize() = %s; want %s", got, expected)
	}
}

func TestError_Deserialize(t *testing.T) {
	e := &Error{}
	input := "-ERR Something went wrong\r\n"
	if err := e.Deserialize(input); err != nil {
		t.Errorf("Deserialize() returned an error: %v", err)
	}
	expectedPrefix := "ERR"
	if e.Prefix != expectedPrefix {
		t.Errorf("Deserialize() Prefix = %s; want %s", e.Prefix, expectedPrefix)
	}
	expectedMessage := "Something went wrong"
	if e.Message != expectedMessage {
		t.Errorf("Deserialize() Message = %s; want %s", e.Message, expectedMessage)
	}
}

func TestInteger_Serialize(t *testing.T) {
	i := &Integer{Value: 42}
	expected := ":42\r\n"
	if got := i.Serialize(); got != expected {
		t.Errorf("Serialize() = %s; want %s", got, expected)
	}
}

func TestInteger_Deserialize(t *testing.T) {
	i := &Integer{}
	input := ":42\r\n"
	if err := i.Deserialize(input); err != nil {
		t.Errorf("Deserialize() returned an error: %v", err)
	}
	expected := 42
	if i.Value != expected {
		t.Errorf("Deserialize() = %d; want %d", i.Value, expected)
	}
}

// Add tests for BulkString and Array here...
func TestBulkString_Serialize(t *testing.T) {
	b := &BulkString{Value: "foobar"}
	expected := "$6\r\nfoobar\r\n"
	if got := b.Serialize(); got != expected {
		t.Errorf("Serialize() = %s; want %s", got, expected)
	}
}

func TestBulkString_Deserialize(t *testing.T) {
	b := &BulkString{}
	input := "$6\r\nfoobar\r\n"
	if err := b.Deserialize(input); err != nil {
		t.Errorf("Deserialize() returned an error: %v", err)
	}
	expected := "foobar"
	if b.Value != expected {
		t.Errorf("Deserialize() = %s; want %s", b.Value, expected)
	}
}

func TestNullBulkString_Serialize(t *testing.T) {
	b := &BulkString{IsNull: true}
	expected := "$-1\r\n"
	if got := b.Serialize(); got != expected {
		t.Errorf("Serialize() = %s; want %s", got, expected)
	}
}

func TestNullBulkString_Deserialize(t *testing.T) {
	b := &BulkString{}
	input := "$-1\r\n"
	if err := b.Deserialize(input); err != nil {
		t.Errorf("Deserialize() returned an error: %v", err)
	}
	if b.IsNull != true {
		t.Errorf("Deserialize() IsNull = %v; want true", b.IsNull)
	}
}

func TestArray_Serialize(t *testing.T) {
	a := &Array{Elements: []Type{
		&SimpleString{Value: "foo"},
		&Error{Prefix: "ERR", Message: "Something went wrong"},
		&Integer{Value: 42},
		&BulkString{Value: "foobar"},
		&Array{Elements: []Type{
			&SimpleString{Value: "foo2"},
		}},
	}}
	expected := "*5\r\n+foo\r\n-ERR Something went wrong\r\n:42\r\n$6\r\nfoobar\r\n*1\r\n+foo2\r\n"
	if got := a.Serialize(); got != expected {
		t.Errorf("Serialize() = %s; want %s", got, expected)
	}
}

func TestArray_Deserialize(t *testing.T) {
	a := &Array{}
	input := "*5\r\n+foo\r\n-ERR Something went wrong\r\n:42\r\n$6\r\nfoobar\r\n*1\r\n+foo2\r\n"
	if err := a.Deserialize(input); err != nil {
		t.Errorf("Deserialize() returned an error: %v", err)
	}
	if len(a.Elements) != 5 {
		t.Errorf("Deserialize() len(Elements) = %d; want 5", len(a.Elements))
	}
	if s, ok := a.Elements[0].(*SimpleString); !ok || s.Value != "foo" {
		t.Errorf("Deserialize() Elements[0] = %v; want SimpleString{Value: \"foo\"}", a.Elements[0])
	}
	if e, ok := a.Elements[1].(*Error); !ok || e.Prefix != "ERR" || e.Message != "Something went wrong" {
		t.Errorf("Deserialize() Elements[1] = %v; want Error{Prefix: \"ERR\", Message: \"Something went wrong\"}", a.Elements[1])
	}
	if i, ok := a.Elements[2].(*Integer); !ok || i.Value != 42 {
		t.Errorf("Deserialize() Elements[2] = %v; want Integer{Value: 42}", a.Elements[2])
	}
	if b, ok := a.Elements[3].(*BulkString); !ok || b.Value != "foobar" {
		t.Errorf("Deserialize() Elements[3] = %v; want BulkString{Value: \"foobar\"}", a.Elements[3])
	}
	if aa, ok := a.Elements[4].(*Array); !ok || len(aa.Elements) != 1 {
		t.Errorf("Deserialize() Elements[4] = %v; want Array{Elements: []Type{SimpleString{Value: \"foo2\"}}}", a.Elements[4])
	}
	if s, ok := a.Elements[4].(*Array).Elements[0].(*SimpleString); !ok || s.Value != "foo2" {
		t.Errorf("Deserialize() Elements[4].Elements[0] = %v; want SimpleString{Value: \"foo2\"}", a.Elements[4].(*Array).Elements[0])
	}
}

func TestNullArray_Serialize(t *testing.T) {
	a := &Array{IsNull: true}
	expected := "*-1\r\n"
	if got := a.Serialize(); got != expected {
		t.Errorf("Serialize() = %s; want %s", got, expected)
	}
}

func TestNullArray_Deserialize(t *testing.T) {
	a := &Array{}
	input := "*-1\r\n"
	if err := a.Deserialize(input); err != nil {
		t.Errorf("Deserialize() returned an error: %v", err)
	}
	if a.IsNull != true {
		t.Errorf("Deserialize() IsNull = %v; want true", a.IsNull)
	}
	if a.Elements != nil {
		t.Errorf("Deserialize() Elements = %v; want nil", a.Elements)
	}
}
