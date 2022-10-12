package dynamo

import (
	"context"
	"encoding/json"
	"fmt"

	"sgorecki.me/golang/event-store/src/internal/es"

	eventstore "sgorecki.me/golang/event-store/src/internal/es/dynamo"
	"sgorecki.me/golang/event-store/src/internal/onlyorder"
)

type Store struct {
	eventStore *eventstore.Store
}

// NewStore creates Store instance.
func NewStore(eventStore *eventstore.Store) *Store {
	return &Store{
		eventStore: eventStore,
	}
}

func (s Store) Load(ctx context.Context, orderID string) (onlyorder.Order, error) {
	streamName := fmt.Sprintf("order-%s", orderID)
	dbEvents, err := s.eventStore.ReadEvents(ctx, streamName)
	if err != nil {
		return onlyorder.Order{}, fmt.Errorf("%v: cannot read events", err)
	}
	if len(dbEvents) == 0 {
		return onlyorder.Order{}, nil // TODO: is it properly handled? how to handle it?
	}

	var order = onlyorder.Order{Entity: es.Entity{Version: eventstore.Version(dbEvents)}}
	for _, event := range dbEvents {
		err := order.When2(event)
		if err != nil {
			return onlyorder.Order{}, err
		}
	}
	return order, nil
}

func (s Store) Save(ctx context.Context, order onlyorder.Order) error {
	if len(order.Changes) == 0 {
		return nil
	}

	dbItems, err := toDBItems(order, order.Changes)
	if err != nil {
		return fmt.Errorf("%v: error converting to DB items", err)
	}

	return s.eventStore.AppendEvents(ctx, dbItems)
}

func toDBItems(order onlyorder.Order, changes []interface{}) ([]eventstore.DBEventItem, error) {
	var items []eventstore.DBEventItem
	for i, change := range changes {
		serializedChange, err := json.Marshal(change)
		if err != nil {
			return nil, err
		}

		var item eventstore.DBEventItem
		switch change.(type) {
		case onlyorder.OrderCreated:
			key := toKey(order, i)
			item = eventstore.DBEventItem{
				EventKey: key,
				Type:     "OrderCreated",
				Data:     string(serializedChange),
			}
		case onlyorder.ItemAdded:
			key := toKey(order, i)
			item = eventstore.DBEventItem{
				EventKey: key,
				Type:     "ItemAdded",
				Data:     string(serializedChange),
			}
		}

		items = append(items, item)
	}
	return items, nil
}

func toKey(order onlyorder.Order, i int) eventstore.EventKey {
	return eventstore.EventKey{
		ID:      fmt.Sprintf("order-%s", order.ID),
		Version: order.Version + i + 1,
	}
}
