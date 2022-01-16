package order

import (
	"time"

	"github.com/rfashwall/kafka-deep-dive/pkg/models"

	"github.com/google/uuid"
)

type OrderConfirmed struct {
	EventBase models.BaseEvent
	EventBody Order
}

func (or OrderConfirmed) ID() uuid.UUID {
	return or.EventBase.EventID
}

func (or OrderConfirmed) Name() string {
	return "OrderConfirmed"
}

func (or OrderConfirmed) Timestamp() time.Time {
	return or.EventBase.EventTimestamp
}

func (or OrderConfirmed) Body() interface{} {
	return or.EventBody
}
