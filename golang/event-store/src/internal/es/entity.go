package es

type whener interface {
	When(interface{})
}

type Applier interface {
	whener
	Apply(interface{})
}

type Entity struct {
	Applier
	changes []interface{}
}

func (e Entity) Apply(event interface{}) {
	e.When(event)
	e.changes = append(e.changes, event)
}
