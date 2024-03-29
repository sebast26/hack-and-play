package dynamo

import (
	"context"
	"encoding/json"
	"fmt"

	"sgorecki.me/golang/event-store/src/internal/es"

	eventstore "sgorecki.me/golang/event-store/src/internal/es/dynamo"
	"sgorecki.me/golang/event-store/src/internal/onlyuser"
)

// UserStore keeps dependencies.
type UserStore struct {
	eventStore *eventstore.Store
}

// NewStore creates UserStore instance.
func NewStore(eventStore *eventstore.Store) *UserStore {
	return &UserStore{
		eventStore: eventStore,
	}
}

func (s UserStore) Load(ctx context.Context, userID string) (onlyuser.User, error) {
	streamName := fmt.Sprintf("user-%s", userID)
	dbEvents, err := s.eventStore.ReadEvents(ctx, streamName)
	if err != nil {
		return onlyuser.User{}, fmt.Errorf("%v: cannot read events", err)
	}
	if len(dbEvents) == 0 {
		return onlyuser.User{}, nil // TODO: is it properly handled? how to handle it?
	}

	var user = onlyuser.User{Entity: es.Entity{Version: eventstore.Version(dbEvents)}}
	for _, event := range dbEvents {
		err := user.When2(event)
		if err != nil {
			return onlyuser.User{}, err
		}
	}
	return user, nil
}

func (s UserStore) Save(ctx context.Context, user onlyuser.User) error {
	if len(user.Changes) == 0 {
		return nil // nothing to do
	}

	dbItems, err := toDBItems(user, user.Changes)
	if err != nil {
		return fmt.Errorf("%v: error converting to DB items", err)
	}

	return s.eventStore.AppendEvents(ctx, dbItems)
}

func toDBItems(user onlyuser.User, changes []interface{}) ([]eventstore.DBEventItem, error) {
	var items []eventstore.DBEventItem
	for i, change := range changes {
		serializedChange, err := json.Marshal(change)
		if err != nil {
			return nil, err
		}

		var item eventstore.DBEventItem
		switch change.(type) {
		case onlyuser.UserCreated:
			key := toKey(user, i)
			item = eventstore.DBEventItem{
				EventKey: key,
				Type:     "UserCreated",
				Data:     string(serializedChange),
			}
		case onlyuser.UserEmailChanged:
			key := toKey(user, i)
			item = eventstore.DBEventItem{
				EventKey: key,
				Type:     "UserEmailChanged",
				Data:     string(serializedChange),
			}
		}

		items = append(items, item)
	}

	return items, nil
}

func toKey(user onlyuser.User, i int) eventstore.EventKey {
	return eventstore.EventKey{
		ID:      fmt.Sprintf("user-%s", user.ID),
		Version: user.Version + i + 1,
	}
}
