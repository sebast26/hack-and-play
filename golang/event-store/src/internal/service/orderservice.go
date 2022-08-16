package service

import (
	"context"

	"sgorecki.me/golang/event-store/src/internal/onlyorder"
	"sgorecki.me/golang/event-store/src/internal/onlyorder/dynamo"
)

type OrderService struct {
	store dynamo.Store
}

func (os OrderService) Handle(command AddOrderItem) {
	ctx := context.Background()
	order, err := os.store.Load(ctx, command.OrderID)
	if err != nil {
		panic("order service load error")
	}

	order.AddItem(command.Item)
	err = os.store.Save(ctx, order)
	if err != nil {
		panic("order service save error")
	}
}

type AddOrderItem struct {
	OrderID string
	Item    onlyorder.OrderItem
}
