package resp

import (
	"fmt"
	"strconv"
	"time"

	"github.com/tn259/cc-redis/database"
)

// https://redis.io/docs/latest/commands/set/
type Set struct {
	key    *BulkString
	value  *BulkString
	expiry *time.Time
}

func NewSet(a *Array) (*Set, error) {
	if len(a.Elements) < 3 {
		return nil, fmt.Errorf("SET command requires at least 3 arguments")
	}
	key := a.Elements[1].(*BulkString)
	value := a.Elements[2].(*BulkString)
	i := 3
	var expiry time.Time
	for i < len(a.Elements) {
		arg := a.Elements[i]
		if s, ok := arg.(*BulkString); ok {
			switch s.Value {
			case "EX":
				if len(a.Elements) < i+1 {
					return nil, fmt.Errorf("SET command requires an argument after EX")
				}
				seconds, err := strconv.Atoi(a.Elements[i+1].(*BulkString).Value)
				if err != nil {
					return nil, fmt.Errorf("SET command requires a valid integer after EX")
				}
				expiry = time.Now().Add(time.Duration(seconds) * time.Second)
				i += 2
			case "PX":
				if len(a.Elements) < i+1 {
					return nil, fmt.Errorf("SET command requires an argument after PX")
				}
				milliseconds, err := strconv.Atoi(a.Elements[i+1].(*BulkString).Value)
				if err != nil {
					return nil, fmt.Errorf("SET command requires a valid integer after PX")
				}
				expiry = time.Now().Add(time.Duration(milliseconds) * time.Millisecond)
				i += 2
			case "EAXT":
				if len(a.Elements) < i+1 {
					return nil, fmt.Errorf("SET command requires an argument after EXAT")
				}
				timestamp, err := strconv.Atoi(a.Elements[i+1].(*BulkString).Value)
				if err != nil {
					return nil, fmt.Errorf("SET command requires a valid integer after EXAT")
				}
				expiry = time.Unix(int64(timestamp), 0)
				i += 2
			case "PXAT":
				if len(a.Elements) < i+1 {
					return nil, fmt.Errorf("SET command requires an argument after PXAT")
				}
				timestamp, err := strconv.Atoi(a.Elements[i+1].(*BulkString).Value)
				if err != nil {
					return nil, fmt.Errorf("SET command requires a valid integer after PXAT")
				}
				expiry = time.Unix(int64(timestamp)/1000, 0)
				i += 2
			default:
				i += 1
			}
			continue
		}
		i += 1
	}
	if expiry.IsZero() {
		return &Set{key: key, value: value, expiry: nil}, nil
	}
	return &Set{key: key, value: value, expiry: &expiry}, nil
}

func (s *Set) Execute() (Type, error) {
	database.Database().Set(s.key.Value, s.value.Value, s.expiry)
	return &SimpleString{Value: "OK"}, nil
}
