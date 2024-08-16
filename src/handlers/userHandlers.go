package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	"github.com/kylerequez/go-user-settings-example/src/middlewares"
	"github.com/kylerequez/go-user-settings-example/src/models"
	"github.com/kylerequez/go-user-settings-example/src/repositories"
	"github.com/kylerequez/go-user-settings-example/src/utils"
	"github.com/kylerequez/go-user-settings-example/src/views"
)

type UserHandler struct {
	server *fiber.App
	ur     *repositories.UserRepository
}

func NewUserHandler(server *fiber.App, ur *repositories.UserRepository) *UserHandler {
	return &UserHandler{
		server: server,
		ur:     ur,
	}
}

func (uh *UserHandler) Init() {
	views := uh.server.Group("")
	views.Get("/", uh.GetHomePage, middlewares.CheckIfSignedIn)
	views.Get("/about-us", uh.GetAboutUsPage, middlewares.CheckIfSignedIn)
	views.Get("/users", uh.GetUsersPage, middlewares.MustBeSignedIn)
	views.Get("/users/create", uh.CreateUser, middlewares.MustBeSignedIn, middlewares.CheckIfAdmin)
	views.Get("/users/:id", uh.GetUserPage, middlewares.MustBeSignedIn)
	views.Get("/my-profile", uh.GetMyProfilePage, middlewares.MustBeSignedIn)

	api := uh.server.Group("/api/v1/users", middlewares.CheckIfSignedIn)
	api.Post("/search", uh.SearchUsers)
	api.Get("/:id/cancel", uh.UserCancel, middlewares.CheckIfAdmin)
	api.Post("", uh.CreateUser, middlewares.CheckIfAdmin)
	api.Get("/:id/update", uh.UpdateUser, middlewares.CheckIfAdmin)
	api.Patch("/:id/update", uh.UpdateUser, middlewares.CheckIfAdmin)
	api.Get("/:id/profile/cancel", uh.UserCancel)
	api.Get("/:id/profile/update", uh.UpdateUser)
	api.Patch("/:id/profile/update", uh.UpdateUser)
	api.Get("/:id/delete", uh.DeleteUser, middlewares.CheckIfAdmin)
	api.Delete("/:id/delete", uh.DeleteUser, middlewares.CheckIfAdmin)
}

func (uh *UserHandler) GetHomePage(c fiber.Ctx) error {
	info := models.AppInfo{
		Title:        "Home",
		CurrentPath:  c.Path(),
		LoggedInUser: RetrieveLoggedInUser(c),
	}

	return Render(c, views.Home(info))
}

func (uh *UserHandler) GetAboutUsPage(c fiber.Ctx) error {
	info := models.AppInfo{
		Title:        "About Us",
		CurrentPath:  c.Path(),
		LoggedInUser: RetrieveLoggedInUser(c),
	}

	return Render(c, views.AboutUs(info))
}

func (uh *UserHandler) GetUsersPage(c fiber.Ctx) error {
	info := models.AppInfo{
		Title:        "Users",
		CurrentPath:  c.Path(),
		LoggedInUser: RetrieveLoggedInUser(c),
	}

	form := models.UsersSearchFormData{}

	users, err := uh.ur.GetAllUsers(c.Context(), info.LoggedInUser.ID)
	if err != nil {
		return nil
	}

	return Render(c, views.Users(info, form, users))
}

func (uh *UserHandler) GetMyProfilePage(c fiber.Ctx) error {
	info := models.AppInfo{
		Title:        "My Profile",
		CurrentPath:  c.Path(),
		LoggedInUser: RetrieveLoggedInUser(c),
	}

	user, err := uh.ur.GetUserById(c.Context(), info.LoggedInUser.ID)
	if err != nil {
		return nil
	}

	return Render(c, views.UserProfile(info, user, nil))
}

func (uh *UserHandler) SearchUsers(c fiber.Ctx) error {
	loggedInUser := RetrieveLoggedInUser(c)

	data, err := ValidateSearchUserForm(c)
	if err != nil {
		return Render(c, views.UsersDisplay(RetrieveLoggedInUser(c), nil))
	}

	users, err := uh.ur.GetUsersByQuery(
		c.Context(),
		loggedInUser.ID,
		data["keyword"],
		data["filter"],
		data["sort"],
	)
	if err != nil {
		return Render(c, views.UsersDisplay(RetrieveLoggedInUser(c), nil))
	}

	return Render(c, views.UsersDisplay(RetrieveLoggedInUser(c), users))
}

func (uh *UserHandler) GetUserPage(c fiber.Ctx) error {
	info := models.AppInfo{
		Title:        "",
		CurrentPath:  c.Path(),
		LoggedInUser: RetrieveLoggedInUser(c),
	}

	id := c.Params("id")
	if id == "" {
		return Redirect(c, "/users")
	}

	oid, err := uuid.Parse(id)
	if err != nil {
		return Redirect(c, "/users")
	}

	user, err := uh.ur.GetUserById(c.Context(), oid)
	if err != nil {
		return Redirect(c, "/users")
	}

	info.Title = fmt.Sprintf("%s's Profile", user.Name)
	return Render(c, views.UserProfile(info, user, nil))
}

func (uh *UserHandler) UserCancel(c fiber.Ctx) error {
	info := models.AppInfo{
		Title:        "",
		CurrentPath:  c.Path(),
		LoggedInUser: RetrieveLoggedInUser(c),
	}

	id := c.Params("id")
	if id == "" {
		return Redirect(c, "/users")
	}

	oid, err := uuid.Parse(id)
	if err != nil {
		return Redirect(c, "/users")
	}

	user, err := uh.ur.GetUserById(c.Context(), oid)
	if err != nil {
		return Redirect(c, "/users")
	}

	info.Title = fmt.Sprintf("%s's Profile", user.Name)
	return Render(c, views.UserProfileContainer(info, user, nil))
}

func (uh *UserHandler) UpdateUser(c fiber.Ctx) error {
	info := models.AppInfo{
		Title:        "",
		CurrentPath:  c.Path(),
		LoggedInUser: RetrieveLoggedInUser(c),
	}

	msg := models.AppMessages{
		Message: nil,
		Errors:  make(map[string]string),
	}

	id := c.Params("id")
	if id == "" {
		return Redirect(c, "/users")
	}

	oid, err := uuid.Parse(id)
	if err != nil {
		return Redirect(c, "/users")
	}

	user, err := uh.ur.GetUserById(c.Context(), oid)
	if err != nil {
		return Redirect(c, "/users")
	}

	if c.Method() == fiber.MethodGet {
		info.Title = fmt.Sprintf("Update %s's Profile", user.Name)
		return Render(c, views.UserUpdateForm(info, user, &msg, nil))
	} else if c.Method() == fiber.MethodPatch {
		form := models.UsersUpdateFormData{}

		if err := ValidateUpdateUserForm(c, user, &msg, &form); err != nil {
			msg.Errors["update-form-err"] = err.Error()
			return Render(c, views.UserUpdateForm(info, user, &msg, &form))
		}

		if err := uh.ur.UpdateUser(c.Context(), *user); err != nil {
			msg.Errors["update-form-err"] = err.Error()
			return Render(c, views.UserUpdateForm(info, user, &msg, &form))
		}

		success := "successfully updated user. please refresh the page for the changes to take effect."
		msg.Message = &success
		return Render(c, views.UserProfileContainer(info, user, &msg))
	}
	return nil
}

func (uh *UserHandler) DeleteUser(c fiber.Ctx) error {
	info := models.AppInfo{
		Title:        "",
		CurrentPath:  c.Path(),
		LoggedInUser: RetrieveLoggedInUser(c),
	}

	msg := models.AppMessages{
		Message: nil,
		Errors:  make(map[string]string),
	}

	id := c.Params("id")
	if id == "" {
		return Redirect(c, "/users")
	}

	oid, err := uuid.Parse(id)
	if err != nil {
		return Redirect(c, "/users")
	}

	user, err := uh.ur.GetUserById(c.Context(), oid)
	if err != nil {
		return Redirect(c, "/users")
	}

	if c.Method() == fiber.MethodGet {
		info.Title = fmt.Sprintf("Delete %s's Profile", user.Name)
		return Render(c, views.UserDeleteForm(info, user, &msg))
	} else if c.Method() == fiber.MethodDelete {
		if err := uh.ur.DeleteUser(c.Context(), oid); err != nil {
			msg.Errors["delete-form-err"] = err.Error()
			return Render(c, views.UserProfileContainer(info, user, &msg))
		}

		return Redirect(c, "/users")
	}
	return nil
}

func (uh *UserHandler) CreateUser(c fiber.Ctx) error {
	info := models.AppInfo{
		Title:        "Create User",
		CurrentPath:  c.Path(),
		LoggedInUser: RetrieveLoggedInUser(c),
	}

	msg := models.AppMessages{
		Message: nil,
		Errors:  make(map[string]string),
	}
	if c.Method() == fiber.MethodGet {
		return Render(c, views.UserCreate(info, nil, nil))
	} else if c.Method() == fiber.MethodPost {
		form := models.CreateUserFormData{}

		user, err := ValidateCreateUserForm(c, &msg, &form)
		if err != nil {
			return Render(c, views.CreateUserForm(&form, &msg))
		}

		password, err := GenerateNewUserPassword(user)
		if err != nil {
			msg.Errors["create-user-form-err"] = err.Error()
			return Render(c, views.CreateUserForm(&form, &msg))
		}

		if err := uh.ur.CreateUser(c.Context(), *user); err != nil {
			msg.Errors["create-user-form-err"] = err.Error()
			return Render(c, views.CreateUserForm(&form, &msg))
		}

		if err := utils.SendEmail(
			"new_account_credential.template.html",
			[]string{user.Email},
			"Go User Settings Example - New User Account",
			struct {
				Email    string
				Password string
			}{
				Email:    user.Email,
				Password: password,
			},
		); err != nil {
			msg.Errors["create-user-form-err"] = err.Error()
			return Render(c, views.CreateUserForm(&form, &msg))
		}

		success := "you have successfully added a new user"
		msg.Message = &success
		return Render(c, views.CreateUserForm(&form, &msg))
	}
	return nil
}
