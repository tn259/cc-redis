package resp

import (
	"fmt"

	"github.com/tn259/cc-redis/database"
)

type lpush struct {
	key    *BulkString
	values []*BulkString
}

func NewLPush(a *Array) (*lpush, error) {
	if len(a.Elements) < 3 {
		return nil, fmt.Errorf("LPUSH command requires at least 3 arguments")
	}
	key := a.Elements[1].(*BulkString)
	values := make([]*BulkString, len(a.Elements)-2)
	for i := 2; i < len(a.Elements); i++ {
		values[i-2] = a.Elements[i].(*BulkString)
	}
	return &lpush{key: key, values: values}, nil
}

func (l *lpush) Execute() (Type, error) {
	db := database.Database()
	for _, value := range l.values {
		err := db.ListLPush(l.key.Value, value.Value)
		if err != nil {
			return nil, err
		}
	}
	return &Integer{Value: len(l.values)}, nil
}
