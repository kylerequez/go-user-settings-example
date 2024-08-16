package models

type RegisterFormData struct {
	Name            string
	Email           string
	Password        string
	ConfirmPassword string
}

type LoginFormData struct {
	Email    string
	Password string
}

type UsersSearchFormData struct {
	Keyword string
	Filter  string
	Sort    string
}

type CreateUserFormData struct {
	Name      string
	Email     string
	Authority string
}

type UsersUpdateFormData struct {
	Name      string
	Email     string
	Authority string
	Theme     string
}
