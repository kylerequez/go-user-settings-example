package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID
	Name      string
	Email     string
	Password  []byte
	Authority string
	Settings  Settings
	CreatedAt time.Time
	UpdatedAt time.Time
}
