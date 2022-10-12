package onlyorder

import (
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	"sgorecki.me/golang/event-store/src/internal/es"
	eventstore "sgorecki.me/golang/event-store/src/internal/es/dynamo"
)

type Order struct {
	es.Entity   // base entity for event sourcing
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

func NewOrderCreated(data string) (OrderCreated, error) {
	var e OrderCreated
	err := json.Unmarshal([]byte(data), &e)
	if err != nil {
		return OrderCreated{}, err
	}
	return e, nil
}

type ItemAdded struct {
	OrderID string
	Item    OrderItem
	Total   int
}

func NewItemAdded(data string) (ItemAdded, error) {
	var e ItemAdded
	err := json.Unmarshal([]byte(data), &e)
	if err != nil {
		return ItemAdded{}, err
	}
	return e, nil
}

func NewOrder() Order {
	order := Order{}
	id := uuid.New()
	order.Apply(&order, OrderCreated{OrderID: id.String()})
	return order
}

func (o *Order) AddItem(item OrderItem) {
	// validateOrderItem
	o.Apply(o, ItemAdded{
		OrderID: o.ID,
		Item:    item,
		Total:   item.Total,
	})
}

func (o *Order) When(event eventstore.DBEventItem) error {
	switch event.Type {
	case "OrderCreated":
		v, err := NewOrderCreated(event.Data)
		if err != nil {
			return err
		}
		o.ID = v.OrderID
	case "ItemAdded":
		v, err := NewItemAdded(event.Data)
		if err != nil {
			return err
		}
		o.OrderItems = append(o.OrderItems, v.Item)
		o.TotalAmount += v.Total
	default:
		return errors.New("unknown event type")
	}
	return nil
}

func (o *Order) When2(event interface{}) {
	switch v := event.(type) {
	case OrderCreated:
		o.ID = v.OrderID
	case ItemAdded:
		o.OrderItems = append(o.OrderItems, v.Item)
		o.TotalAmount += v.Total
	default:
		panic("unknown event type!")
	}
}
