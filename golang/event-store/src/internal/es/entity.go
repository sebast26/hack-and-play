package es

//type whener interface {
//	When(interface{})
//}
//
//type Applier interface {
//	whener
//	Apply(interface{})
//}
//
//type Entity struct {
//	Applier
//	changes []interface{}
//}
//
//func (e Entity) Apply(event interface{}) {
//	e.When(event)
//	e.changes = append(e.changes, event)
//}

type Entity struct {
	ID      string
	Version int
	Changes []interface{}
}

// Apply function is available to all structs that have Entity embedded
func (e Entity) Apply(event interface{}) {
	// TODO: e.When(event)
	e.Changes = append(e.Changes, event)
}

type UserEntity struct {
	Entity
	Name  string
	Email string
}

type OrderEntity struct {
	Entity
	Items []string
	Total int
}

func bothEntitiesHaveApplyAndCommonFields() {
	user := UserEntity{}
	order := OrderEntity{}

	user.Apply("some event here")
	order.Apply("some event here")

	_ = user.Changes
	_ = order.Changes
	_ = user.Version
	_ = order.Version
}
