package handlers

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/storage/postgres/v3"

	"github.com/kylerequez/go-user-settings-example/src/middlewares"
	"github.com/kylerequez/go-user-settings-example/src/models"
	"github.com/kylerequez/go-user-settings-example/src/repositories"
	"github.com/kylerequez/go-user-settings-example/src/utils"
	"github.com/kylerequez/go-user-settings-example/src/views"
)

type AuthHandler struct {
	server *fiber.App
	ur     *repositories.UserRepository
	store  *postgres.Storage
}

func NewAuthHandler(
	server *fiber.App,
	ur *repositories.UserRepository,
	store *postgres.Storage,
) *AuthHandler {
	return &AuthHandler{
		server: server,
		ur:     ur,
		store:  store,
	}
}

func (ah *AuthHandler) Init() {
	api := ah.server.Group("/api/v1/auth")
	api.Post("/login", ah.LoginUser)
	api.Post("/register", ah.RegisterUser)
	api.Get("/logout", ah.LogoutUser, middlewares.MustBeSignedIn)

	views := ah.server.Group("")
	views.Get("/login", ah.GetLoginPage)
	views.Get("/register", ah.GetRegisterPage)
}

func (ah *AuthHandler) LoginUser(c fiber.Ctx) error {
	msg := models.AppMessages{
		Message: nil,
		Errors:  make(map[string]string),
	}

	form := models.LoginFormData{}

	user, err := ValidateLoginForm(c, &msg, &form)
	if err != nil {
		msg.Errors["login-form-err"] = err.Error()
		return Render(c, views.LoginForm(&form, &msg))
	}

	foundUser, err := ah.ur.GetUserByEmail(c.Context(), user.Email)
	if err != nil {
		msg.Errors["login-form-err"] = err.Error()
		return Render(c, views.LoginForm(&form, &msg))
	}

	err = utils.VerifyPassword(foundUser.Password, user.Password)
	if err != nil {
		msg.Errors["login-form-err"] = err.Error()
		return Render(c, views.LoginForm(&form, &msg))
	}

	token, err := utils.GenerateJwtToken(*foundUser)
	if err != nil {
		msg.Errors["login-form-err"] = err.Error()
		return Render(c, views.LoginForm(&form, &msg))
	}

	err = StoreUserSession(c, ah.store, foundUser.ID.String(), token)
	if err != nil {
		msg.Errors["login-form-err"] = err.Error()
		return Render(c, views.LoginForm(&form, &msg))
	}

	return Redirect(c, "/users")
}

func (ah *AuthHandler) LogoutUser(c fiber.Ctx) error {
	loggedInUser := RetrieveLoggedInUser(c)
	if loggedInUser == nil {
		return Redirect(c, "/login")
	}

	if err := ah.store.Delete(loggedInUser.ID.String()); err != nil {
		return Redirect(c, "/login")
	}

	return Redirect(c, "/login")
}

func (ah *AuthHandler) RegisterUser(c fiber.Ctx) error {
	msg := models.AppMessages{
		Message: nil,
		Errors:  make(map[string]string),
	}

	form := models.RegisterFormData{}

	user, err := ValidateRegisterForm(c, &msg, &form)
	if err != nil {
		msg.Errors["reg-form-err"] = err.Error()
		return Render(c, views.RegisterForm(&form, &msg))
	}

	err = ah.ur.CreateUser(c.Context(), *user)
	if err != nil {
		msg.Errors["reg-form-err"] = err.Error()
		return Render(c, views.RegisterForm(&form, &msg))
	}

	return Redirect(c, "/login")
}

func (ah *AuthHandler) GetRegisterPage(c fiber.Ctx) error {
	c.ClearCookie("go-user-settings-example-cookie")

	info := models.AppInfo{
		Title:       "Register",
		CurrentPath: c.Path(),
	}

	return Render(c, views.Register(info, nil, nil))
}

func (ah *AuthHandler) GetLoginPage(c fiber.Ctx) error {
	c.ClearCookie("go-user-settings-example-cookie")

	info := models.AppInfo{
		Title:       "Login",
		CurrentPath: c.Path(),
	}

	return Render(c, views.Login(info, nil, nil))
}
