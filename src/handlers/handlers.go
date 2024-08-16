package handlers

import (
	"context"

	"github.com/gofiber/fiber/v3"

	"github.com/kylerequez/go-user-settings-example/src/db"
	"github.com/kylerequez/go-user-settings-example/src/repositories"
)

func Init(server *fiber.App, ctx context.Context) error {
	if err := db.Connect(ctx); err != nil {
		return err
	}

	if err := db.Ping(ctx); err != nil {
		return err
	}

	store, err := db.InitStore()
	if err != nil {
		return err
	}

	ur := repositories.NewUserRepository(db.DB, "users")
	uh := NewUserHandler(server, ur)
	uh.Init()

	ah := NewAuthHandler(server, ur, store)
	ah.Init()

	return nil
}
