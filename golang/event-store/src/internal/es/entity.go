package es

type EventSourcer interface {
	Apply(EventSourcer, interface{})
	When()
}

type Entity struct {
	ID      string
	Version int
	Changes []interface{}
}

func (e *Entity) Apply(entity EventSourcer, event interface{}) {
	entity.When()
	e.Changes = append(e.Changes, event)
}
