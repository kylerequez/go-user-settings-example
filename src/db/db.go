package db

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	"github.com/kylerequez/go-user-settings-example/src/utils"
)

var DB *pgx.Conn

func Connect(ctx context.Context) error {
	uri, err := utils.GetEnv("DB_URI")
	if err != nil {
		return err
	}

	if uri == nil {
		return errors.New("")
	}

	conn, err := pgx.Connect(ctx, *uri)
	if err != nil {
		return err
	}

	DB = conn
	return nil
}

func Close(ctx context.Context) error {
	if err := DB.Close(ctx); err != nil {
		return err
	}
	return nil
}

func Ping(ctx context.Context) error {
	if err := DB.Ping(ctx); err != nil {
		return err
	}
	return nil
}
