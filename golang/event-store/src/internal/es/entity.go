package es

// EventSourcer interface for event-based entities (following Event Sourcing pattern).
type EventSourcer interface {
	Apply(EventSourcer, interface{})
	When()
}

// Entity is a base type for all event-based entities.
type Entity struct {
	ID      string
	Version int
	Changes []interface{}
}

// Apply generic method that applies event on given entity.
func (e *Entity) Apply(entity EventSourcer, event interface{}) {
	entity.When()
	e.Changes = append(e.Changes, event)
}
