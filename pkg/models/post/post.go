package post

import "github.com/google/uuid"

type Post struct {
	ID uuid.UUID `json:"id,omitempty"`
}
