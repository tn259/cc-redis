package resp

import (
	"fmt"

	"github.com/tn259/cc-redis/database"
)

type Delete struct {
	keys []*BulkString
}

func NewDelete(a *Array) (*Delete, error) {
	if len(a.Elements) < 2 {
		return nil, fmt.Errorf("DELETE command requires at least 2 arguments")
	}
	keys := make([]*BulkString, len(a.Elements)-1)
	for i := 1; i < len(a.Elements); i++ {
		keys[i-1] = a.Elements[i].(*BulkString)
	}
	return &Delete{keys: keys}, nil
}

func (e *Delete) Execute() (Type, error) {
	db := database.Database()
	for _, key := range e.keys {
		db.Delete(key.Value)
	}
	return &Integer{Value: len(e.keys)}, nil
}
