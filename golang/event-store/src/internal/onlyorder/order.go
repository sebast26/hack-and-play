package onlyorder

import "github.com/google/uuid"

type Order struct {
	OrderID     string
	OrderItems  []OrderItem
	TotalAmount int

	// TODO: this should be generic! belongs to Entity class
	Changes []interface{}
	Version int
}

type OrderItem struct {
	Name  string
	Total int
}

type OrderCreated struct {
	OrderID string
}

type ItemAdded struct {
	OrderID string
	Item    OrderItem
	Total   int
}

func NewOrder() Order {
	order := Order{}
	id := uuid.New()
	order.Apply(OrderCreated{OrderID: id.String()})
	return order
}

func (o *Order) AddItem(item OrderItem) {
	o.Apply(ItemAdded{
		OrderID: o.OrderID,
		Item:    item,
		Total:   item.Total,
	})
}

// TODO: this should be generic! belongs to Entity class
func (o *Order) Apply(event interface{}) {
	o.When(event)
	o.Changes = append(o.Changes, event)
}

// TODO: this should be abstract method reimplemented in all subclasses!
func (o *Order) When(event interface{}) {
	switch v := event.(type) {
	case OrderCreated:
		o.OrderID = v.OrderID
	case ItemAdded:
		o.OrderItems = append(o.OrderItems, v.Item)
		o.TotalAmount += v.Total
	default:
		panic("unknown event type!")
	}
	o.Version++
}
