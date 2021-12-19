package model

import "github.com/google/uuid"

type Post struct {
	ID uuid.UUID `json:"id,omitempty"`
}
