package dynamo

import (
	"context"

	"sgorecki.me/golang/event-store/src/internal/onlyorder"
)

type Store struct {
}

func (s Store) Load(ctx context.Context, orderID string) (onlyorder.Order, error) {
	return onlyorder.Order{}, nil
}

func (s Store) Save(ctx context.Context, order onlyorder.Order) error {
	return nil
}
