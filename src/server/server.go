package server

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/static"

	"github.com/kylerequez/go-user-settings-example/src/handlers"
	"github.com/kylerequez/go-user-settings-example/src/utils"
)

func Init() {
	if err := utils.LoadEnvVariables(); err != nil {
		panic(err)
	}

	server := fiber.New(fiber.Config{
		AppName: "go-user-settings-example",
	})

	server.Use("/public", static.New("./src/public"))
	server.Use(logger.New(logger.ConfigDefault))

	if err := handlers.Init(server, context.Background()); err != nil {
		panic(err)
	}

	port, err := utils.GetEnv("SERVER_PORT")
	if err != nil {
		panic(err)
	}
	server.Listen(fmt.Sprintf(":%s", *port))
}
