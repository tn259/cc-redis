package resp

import (
	"fmt"
	"strconv"

	"github.com/tn259/cc-redis/database"
)

type Decr struct {
	key *BulkString
}

func (d *Decr) Execute() (Type, error) {
	db := database.Database()
	value, ok := db.Get(d.key.Value)
	if !ok {
		value = "0"
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return nil, fmt.Errorf("DECR key %s has value %s which is not an integer", d.key.Value, value)
	}
	intValue--
	db.Set(d.key.Value, strconv.Itoa(intValue), nil)
	return &Integer{Value: intValue}, nil
}
