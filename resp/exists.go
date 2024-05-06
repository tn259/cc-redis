package resp

import (
	"fmt"

	"github.com/tn259/cc-redis/database"
)

type Exists struct {
	keys []*BulkString
}

func NewExists(a *Array) (*Exists, error) {
	if len(a.Elements) < 2 {
		return nil, fmt.Errorf("EXISTS command requires at least 2 arguments")
	}
	keys := make([]*BulkString, len(a.Elements)-1)
	for i := 1; i < len(a.Elements); i++ {
		keys[i-1] = a.Elements[i].(*BulkString)
	}
	return &Exists{keys: keys}, nil
}

func (e *Exists) Execute() (Type, error) {
	exists := 0
	db := database.Database()
	for _, key := range e.keys {
		if _, ok := db.Get(key.Value); ok {
			exists++
		}
	}
	return &Integer{Value: exists}, nil
}
