package es

import eventstore "sgorecki.me/golang/event-store/src/internal/es/dynamo"

// EventSourcer interface for event-based entities (following Event Sourcing pattern).
type EventSourcer interface {
	Apply(EventSourcer, eventstore.DBEventItem) error
	When(event eventstore.DBEventItem) error
}

// Entity is a base type for all event-based entities.
type Entity struct {
	ID      string
	Version int
	Changes []interface{}
}

// Apply generic method that applies event on given entity.
func (e *Entity) Apply(entity EventSourcer, event eventstore.DBEventItem) error {
	err := entity.When(event)
	if err != nil {
		return err
	}
	e.Changes = append(e.Changes, event)
	return nil
}
