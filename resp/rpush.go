package resp

import (
	"fmt"

	"github.com/tn259/cc-redis/database"
)

type rpush struct {
	key    *BulkString
	values []*BulkString
}

func NewRPush(a *Array) (*rpush, error) {
	if len(a.Elements) < 3 {
		return nil, fmt.Errorf("RPUSH command requires at least 3 arguments")
	}
	key := a.Elements[1].(*BulkString)
	values := make([]*BulkString, len(a.Elements)-2)
	for i := 2; i < len(a.Elements); i++ {
		values[i-2] = a.Elements[i].(*BulkString)
	}
	return &rpush{key: key, values: values}, nil
}

func (r *rpush) Execute() (Type, error) {
	db := database.Database()
	for _, value := range r.values {
		err := db.ListRPush(r.key.Value, value.Value)
		if err != nil {
			return nil, err
		}
	}
	return &Integer{Value: len(r.values)}, nil
}
