package handlers

import (
	"errors"
	"net/mail"
	"strings"
	"time"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/adaptor"
	"github.com/gofiber/storage/postgres/v3"

	"github.com/kylerequez/go-user-settings-example/src/models"
	"github.com/kylerequez/go-user-settings-example/src/repositories"
	"github.com/kylerequez/go-user-settings-example/src/utils"
)

func Render(
	c fiber.Ctx,
	component templ.Component,
	options ...func(*templ.ComponentHandler),
) error {
	componentHandler := templ.Handler(component)
	for _, o := range options {
		o(componentHandler)
	}
	return adaptor.HTTPHandler(componentHandler)(c)
}

func Redirect(
	c fiber.Ctx,
	path string,
) error {
	c.Set("HX-Redirect", path)
	return c.SendStatus(fiber.StatusOK)
}

type ErrorResponse struct {
	Status int
	Err    error
}

func ValidateUser(c fiber.Ctx) (*models.User, ErrorResponse) {
	type RequestBody struct {
		Name            string `json:"name"`
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}

	body := new(RequestBody)
	if err := c.Bind().Body(body); err != nil {
		return nil, ErrorResponse{
			Status: fiber.StatusUnprocessableEntity,
			Err:    err,
		}
	}

	if strings.TrimSpace(body.Name) == "" {
		return nil, ErrorResponse{
			Status: fiber.StatusBadRequest,
			Err:    errors.New("name must not be empty"),
		}
	}

	if strings.TrimSpace(body.Email) == "" {
		return nil, ErrorResponse{
			Status: fiber.StatusBadRequest,
			Err:    errors.New("email must not be empty"),
		}
	}

	if _, err := mail.ParseAddress(body.Email); err != nil {
		return nil, ErrorResponse{
			Status: fiber.StatusBadRequest,
			Err:    errors.New("email is invalid"),
		}
	}

	if strings.TrimSpace(body.Password) == "" {
		return nil, ErrorResponse{
			Status: fiber.StatusBadRequest,
			Err:    errors.New("password must not be empty"),
		}
	}

	if strings.TrimSpace(body.ConfirmPassword) == "" {
		return nil, ErrorResponse{
			Status: fiber.StatusBadRequest,
			Err:    errors.New("confirm password must not be empty"),
		}
	}

	if body.Password != body.ConfirmPassword {
		return nil, ErrorResponse{
			Status: fiber.StatusUnprocessableEntity,
			Err:    errors.New("passwords much match"),
		}
	}

	hashedPassword, err := utils.HashPassword(body.Password)
	if err != nil {
		return nil, ErrorResponse{
			Status: fiber.StatusInternalServerError,
			Err:    err,
		}
	}

	currentTime := time.Now()

	user := &models.User{
		Name:      body.Name,
		Email:     body.Email,
		Password:  hashedPassword,
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
	}

	return user, ErrorResponse{
		Status: fiber.StatusOK,
		Err:    nil,
	}
}

func RetrieveLoggedInUser(c fiber.Ctx) *models.User {
	user := c.Locals("loggedInUser")
	loggedInUser, ok := user.(*models.User)
	if !ok {
		return nil
	}
	return loggedInUser
}

func ValidateLoginForm(
	c fiber.Ctx,
	msg *models.AppMessages,
	form *models.LoginFormData,
) (*models.User, error) {
	type RequestBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	body := new(RequestBody)
	if err := c.Bind().Body(body); err != nil {
		return nil, err
	}

	hasError := false
	if strings.TrimSpace(body.Email) == "" {
		msg.Errors["login-form-email-err"] = "email must not be empty"
		hasError = true
	}

	if _, err := mail.ParseAddress(body.Email); err != nil {
		msg.Errors["login-form-email-err"] = "email is invalid: " + err.Error()
		hasError = true
	}

	if strings.TrimSpace(body.Password) == "" {
		msg.Errors["login-form-password-err"] = "password must not be empty"
		hasError = true
	}

	if len(strings.TrimSpace(body.Password)) < 12 {
		msg.Errors["login-form-password-err"] = "password must not be less than 12 characters"
		hasError = true
	}

	if hasError {
		return nil, errors.New("")
	}

	user := &models.User{
		Email:    body.Email,
		Password: []byte(body.Password),
	}

	return user, nil
}

func StoreUserSession(c fiber.Ctx, store *postgres.Storage, id string, token string) error {
	cookie := new(fiber.Cookie)
	cookie.Name = "go-user-settings-example-cookie"
	cookie.Value = token
	exp := time.Now().Add(24 * time.Hour)
	cookie.Expires = exp
	c.Cookie(cookie)

	store.Set(id, []byte(token), exp.Sub(time.Now()))

	return nil
}

func ValidateRegisterForm(
	c fiber.Ctx,
	msg *models.AppMessages,
	form *models.RegisterFormData,
) (*models.User, error) {
	type RequestBody struct {
		Name            string `json:"name"`
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}

	body := new(RequestBody)
	if err := c.Bind().Body(body); err != nil {
		return nil, err
	}

	hasError := false
	if strings.TrimSpace(body.Name) == "" {
		msg.Errors["reg-form-name-err"] = "name must not be empty"
		hasError = true
	}

	if strings.TrimSpace(body.Email) == "" {
		msg.Errors["reg-form-email-err"] = "email must not be empty"
		hasError = true
	}

	if _, err := mail.ParseAddress(body.Email); err != nil {
		msg.Errors["reg-form-email-err"] = "email is invalid: " + err.Error()
		hasError = true
	}

	if strings.TrimSpace(body.Password) == "" {
		msg.Errors["reg-form-password-err"] = "password must not be empty"
		hasError = true
	}

	if len(strings.TrimSpace(body.Password)) < 12 {
		msg.Errors["reg-form-password-err"] = "password must not be less than 12 characters"
		hasError = true
	}

	if strings.TrimSpace(body.ConfirmPassword) == "" {
		msg.Errors["reg-form-confirm-password-err"] = "confirm password must not be empty"
		hasError = true
	}

	if len(strings.TrimSpace(body.ConfirmPassword)) < 12 {
		msg.Errors["reg-form-confirm-password-err"] = "confirm password must not be less than 12 characters"
		hasError = true
	}

	if body.Password != body.ConfirmPassword {
		msg.Errors["reg-form-password-err"] = "passswords must match"
		msg.Errors["reg-form-confirm-password-err"] = "passswords must match"
		hasError = true
	}

	form.Name = body.Name
	form.Email = body.Email
	form.Password = body.Password
	form.ConfirmPassword = body.ConfirmPassword

	if hasError {
		return nil, errors.New("")
	}

	hashedPassword, err := utils.HashPassword(body.Password)
	if err != nil {
		return nil, err
	}

	currentTime := time.Now()

	user := &models.User{
		Name:      body.Name,
		Email:     body.Email,
		Authority: repositories.AUTH_USER,
		Password:  hashedPassword,
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
	}

	form = nil
	return user, nil
}

func ValidateSearchUserForm(
	c fiber.Ctx,
) (map[string]*string, error) {
	type RequestBody struct {
		Keyword string `form:"keyword" json:"keyword" xml:"keyword"`
		Filter  string `form:"filter"  json:"filter"  xml:"filter"`
		Sort    string `form:"sort"    json:"sort"    xml:"sort"`
	}

	data := map[string]*string{
		"keyword": nil,
		"filter":  nil,
		"sort":    nil,
	}

	body := new(RequestBody)
	if err := c.Bind().Body(body); err != nil {
		return nil, err
	}

	if body.Keyword != "" {
		data["keyword"] = &body.Keyword
	}

	if body.Filter != "" {
		data["filter"] = &body.Filter
	}

	if body.Sort != "" {
		data["sort"] = &body.Sort
	}

	return data, nil
}

func ValidateUpdateUserForm(
	c fiber.Ctx,
	user *models.User,
	msg *models.AppMessages,
	form *models.UsersUpdateFormData,
) error {
	type RequestBody struct {
		Name      string `json:"name"`
		Email     string `json:"email"`
		Authority string `json:"authority"`
		Theme     string `json:"theme"`
	}

	body := new(RequestBody)
	if err := c.Bind().Body(body); err != nil {
		return err
	}

	hasError := false
	if strings.TrimSpace(body.Name) == "" {
		msg.Errors["update-form-name-err"] = "name must not be empty"
		hasError = true
	}

	if strings.TrimSpace(body.Email) == "" {
		msg.Errors["update-form-email-err"] = "email must not be empty"
		hasError = true
	}

	if _, err := mail.ParseAddress(body.Email); err != nil {
		msg.Errors["update-form-email-err"] = "email is invalid: " + err.Error()
		hasError = true
	}

	if strings.TrimSpace(body.Authority) == "" {
		msg.Errors["update-form-auth-err"] = "authority must not be empty"
		hasError = true
	}

	if strings.TrimSpace(body.Theme) == "" {
		msg.Errors["update-form-theme-err"] = "theme must not be empty"
		hasError = true
	}

	form.Name = body.Name
	form.Email = body.Email
	form.Authority = body.Authority
	form.Theme = body.Theme

	if hasError {
		return errors.New("")
	}

	currentTime := time.Now()

	user.Name = body.Name
	user.Email = body.Email
	user.Authority = body.Authority
	user.Settings.Theme = body.Theme
	user.UpdatedAt = currentTime

	form = nil
	return nil
}

func ValidateCreateUserForm(
	c fiber.Ctx,
	msg *models.AppMessages,
	form *models.CreateUserFormData,
) (*models.User, error) {
	type RequestBody struct {
		Name      string `json:"name"`
		Email     string `json:"email"`
		Authority string `json:"authority"`
	}

	body := new(RequestBody)
	if err := c.Bind().Body(body); err != nil {
		return nil, err
	}

	hasError := false
	if strings.TrimSpace(body.Name) == "" {
		msg.Errors["create-user-form-name-err"] = "name must not be empty"
		hasError = true
	}

	if strings.TrimSpace(body.Email) == "" {
		msg.Errors["create-user-form-email-err"] = "email must not be empty"
		hasError = true
	}

	if _, err := mail.ParseAddress(body.Email); err != nil {
		msg.Errors["create-user-form-email-err"] = "email is invalid: " + err.Error()
		hasError = true
	}

	if strings.TrimSpace(body.Authority) == "" {
		msg.Errors["create-user-form-authority-err"] = "authority must not be empty"
		hasError = true
	}

	form.Name = body.Name
	form.Email = body.Email
	form.Authority = body.Authority

	if hasError {
		return nil, errors.New("")
	}

	currentTime := time.Now()

	user := &models.User{
		Name:      body.Name,
		Email:     body.Email,
		Authority: body.Authority,
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
	}

	form = nil
	return user, nil
}

func GenerateNewUserPassword(user *models.User) (string, error) {
	generatedPassword := utils.GenerateRandomPassword(12)
	hashedPassword, err := utils.HashPassword(generatedPassword)
	if err != nil {
		return "", err
	}

	user.Password = hashedPassword
	return generatedPassword, err
}
