package resp

import (
	"fmt"

	"github.com/tn259/cc-redis/database"
)

type lrange struct {
	key   *BulkString
	start *BulkString
	stop  *BulkString
}

func NewLRange(a *Array) (*lrange, error) {
	if len(a.Elements) != 4 {
		return nil, fmt.Errorf("LPUSH command requires 3 arguments")
	}
	key := a.Elements[1].(*BulkString)
	start := a.Elements[2].(*BulkString)
	stop := a.Elements[3].(*BulkString)
	return &lrange{key: key, start: start, stop: stop}, nil
}

func (l *lrange) Execute() (Type, error) {
	db := database.Database()
	values, err := db.ListRange(l.key.Value, l.start.Value, l.stop.Value)
	if err != nil {
		return nil, err
	}
	elements := make([]Type, len(values)) // Convert elements slice to []Type
	for i, value := range values {
		elements[i] = &BulkString{Value: value}
	}
	return &Array{Elements: elements}, nil
}
