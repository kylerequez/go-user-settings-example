package db

import (
	"github.com/gofiber/storage/postgres/v3"

	"github.com/kylerequez/go-user-settings-example/src/utils"
)

func InitStore() (*postgres.Storage, error) {
	uri, err := utils.GetEnv("DB_URI")
	if err != nil {
		return nil, err
	}

	store := postgres.New(postgres.Config{
		ConnectionURI: *uri,
	})

	return store, nil
}
