package order

import (
	"time"

	"github.com/rfashwall/kafka-deep-dive/pkg/models"

	"github.com/google/uuid"
)

type OrderReceived struct {
	EventBase models.BaseEvent
	EventBody Order
}

func (or OrderReceived) ID() uuid.UUID {
	return or.EventBase.EventID
}

func (or OrderReceived) Name() string {
	return "OrderReceived"
}

func (or OrderReceived) Timestamp() time.Time {
	return or.EventBase.EventTimestamp
}

func (or OrderReceived) Body() interface{} {
	return or.EventBody
}
