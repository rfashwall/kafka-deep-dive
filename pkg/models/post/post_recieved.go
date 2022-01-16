package post

import (
	"pkg/pkg/models"
	"time"

	"github.com/google/uuid"
)

type PostReceived struct {
	EventBase models.BaseEvent
	EventBody Post
}

func (or PostReceived) ID() uuid.UUID {
	return or.EventBase.EventID
}

func (or PostReceived) Name() string {
	return "PostReceived"
}

func (or PostReceived) Timestamp() time.Time {
	return or.EventBase.EventTimestamp
}

func (or PostReceived) Body() interface{} {
	return or.EventBody
}
