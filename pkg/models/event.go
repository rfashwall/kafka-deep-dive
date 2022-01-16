package models

import (
	"time"

	"github.com/google/uuid"
)

type Event interface {
	ID() uuid.UUID
	Name() string
	Timestamp() time.Time
	Body() interface{}
}

// BaseEvent represents common properties of an event
type BaseEvent struct {
	EventID        uuid.UUID
	EventTimestamp time.Time
}
