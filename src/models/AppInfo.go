package models

type AppInfo struct {
	Title        string
	CurrentPath  string
	LoggedInUser *User
}
