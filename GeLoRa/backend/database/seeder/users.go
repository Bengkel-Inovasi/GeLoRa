package databaseseeder

import (
	_ "embed"
	"encoding/json"
)

//go:embed users.json
var usersJSON []byte

type UserSeedFile struct {
	Seed []struct {
		Id       int64  `json:"id"`
		Name     string `json:"name"`
		Username string `json:"username"`
		Password string `json:"password"`
		Role     string `json:"role"`
		Bio      string `json:"bio"`
	} `json:"seed"`
}

func LoadUserSeedFile() (UserSeedFile, error) {
	var seed UserSeedFile
	if err := json.Unmarshal(usersJSON, &seed); err != nil {
		return UserSeedFile{}, err
	}
	return seed, nil
}
