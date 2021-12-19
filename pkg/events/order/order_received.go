package order

import (
	"pkg/pkg/events"
	"pkg/pkg/order/pkg/model"
	"time"

	"github.com/google/uuid"
)

type OrderReceived struct {
	EventBase events.BaseEvent
	EventBody model.Order
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
