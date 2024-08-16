package middlewares

import (
	"github.com/gofiber/fiber/v3"

	"github.com/kylerequez/go-user-settings-example/src/models"
	"github.com/kylerequez/go-user-settings-example/src/repositories"
)

func CheckIfAdmin(c fiber.Ctx) error {
	user := c.Locals("loggedInUser")
	loggedInUser, ok := user.(*models.User)
	if !ok {
		return c.Redirect().To(c.Request().URI().String())
	}

	if loggedInUser.Authority != repositories.AUTH_ADMIN {
		return c.Redirect().To(c.Request().URI().String())
	}

	return c.Next()
}
