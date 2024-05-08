package resp

import "github.com/tn259/cc-redis/database"

type Save struct {
}

func (s *Save) Execute() (Type, error) {
	db := database.Database()
	err := db.Save()
	if err != nil {
		return nil, err
	}
	return &SimpleString{Value: "OK"}, nil
}
