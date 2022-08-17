package onlyorder

import (
	"github.com/google/uuid"
	"sgorecki.me/golang/event-store/src/internal/es"
)

type Order struct {
	es.Entity
	OrderItems  []OrderItem
	TotalAmount int
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
	order.Apply(&order, OrderCreated{OrderID: id.String()})
	return order
}

func (o *Order) AddItem(item OrderItem) {
	o.Apply(o, ItemAdded{
		OrderID: o.ID,
		Item:    item,
		Total:   item.Total,
	})
}

func (o *Order) When(event interface{}) {
	switch v := event.(type) {
	case OrderCreated:
		o.ID = v.OrderID
	case ItemAdded:
		o.OrderItems = append(o.OrderItems, v.Item)
		o.TotalAmount += v.Total
	default:
		panic("unknown event type!")
	}
	o.Version++
}
