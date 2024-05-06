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
		db.Set(d.key.Value, "0", nil)
		return &Integer{Value: 0}, nil
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return nil, fmt.Errorf("DECR key %s has value %s which is not an integer", d.key.Value, value)
	}
	intValue--
	db.Set(d.key.Value, strconv.Itoa(intValue), nil)
	return &Integer{Value: intValue}, nil
}
