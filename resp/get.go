package resp

import "github.com/tn259/cc-redis/database"

// https://redis.io/docs/latest/commands/get/
type Get struct {
	key *BulkString
}

func (g *Get) Execute() (Type, error) {
	value, ok := database.Database().Get(g.key.Value)
	if !ok {
		return &BulkString{IsNull: true}, nil
	}
	return &BulkString{Value: value}, nil
}
