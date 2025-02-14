package resp

import (
	"fmt"
	"strconv"

	"github.com/tn259/cc-redis/database"
)

type Incr struct {
	key *BulkString
}

func (i *Incr) Execute() (Type, error) {
	db := database.Database()
	value, ok := db.Get(i.key.Value)
	if !ok {
		value = "0"
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return nil, fmt.Errorf("INCR key %s has value %s which is not an integer", i.key.Value, value)
	}
	intValue++
	db.Set(i.key.Value, strconv.Itoa(intValue), nil)
	return &Integer{Value: intValue}, nil
}
