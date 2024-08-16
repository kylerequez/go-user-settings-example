package middlewares

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	"github.com/kylerequez/go-user-settings-example/src/db"
	"github.com/kylerequez/go-user-settings-example/src/repositories"
	"github.com/kylerequez/go-user-settings-example/src/utils"
)

func MustBeSignedIn(c fiber.Ctx) error {
	cookie := c.Cookies("go-user-settings-example-cookie", "")

	if cookie == "" {
		c.ClearCookie("go-user-settings-example-cookie")
		c.Set("hx-redirect", "/login")
		return c.Redirect().To("/login")
	}

	claims, err := utils.ExtractClaims(cookie)
	if err != nil {
		c.ClearCookie("go-user-settings-example-cookie")
		c.Set("hx-redirect", "/login")
		return c.Redirect().To("/login")
	}

	oid, err := uuid.Parse(claims["ID"].(string))
	if err != nil {
		c.ClearCookie("go-user-settings-example-cookie")
		c.Set("hx-redirect", "/login")
		return c.Redirect().To("/login")
	}

	ur := repositories.NewUserRepository(db.DB, "users")
	user, err := ur.GetUserById(c.Context(), oid)
	if err != nil {
		c.ClearCookie("go-user-settings-example-cookie")
		c.Set("hx-redirect", "/login")
		return c.Redirect().To("/login")
	}

	c.Locals("loggedInUser", user)
	return c.Next()
}

func CheckIfSignedIn(c fiber.Ctx) error {
	cookie := c.Cookies("go-user-settings-example-cookie", "")

	if cookie == "" {
		c.ClearCookie("go-user-settings-example-cookie")
		return c.Next()
	}

	claims, err := utils.ExtractClaims(cookie)
	if err != nil {
		c.ClearCookie("go-user-settings-example-cookie")
		c.Set("hx-redirect", "/login")
		return c.Redirect().To("/login")
	}

	oid, err := uuid.Parse(claims["ID"].(string))
	if err != nil {
		c.ClearCookie("go-user-settings-example-cookie")
		c.Set("hx-redirect", "/login")
		return c.Redirect().To("/login")
	}

	ur := repositories.NewUserRepository(db.DB, "users")
	user, err := ur.GetUserById(c.Context(), oid)
	if err != nil {
		c.ClearCookie("go-user-settings-example-cookie")
		c.Set("hx-redirect", "/login")
		return c.Redirect().To("/login")
	}

	c.Locals("loggedInUser", user)
	return c.Next()
}
