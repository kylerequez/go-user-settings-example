package models

import "github.com/google/uuid"

type Settings struct {
	ID     uuid.UUID
	UserID uuid.UUID
	Theme  string
}
