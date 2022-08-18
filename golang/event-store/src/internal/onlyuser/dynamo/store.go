package dynamo

import (
	"context"
	"encoding/json"
	"fmt"

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

	events, err := loadEvents(dbEvents)
	if err != nil {
		return onlyuser.User{}, fmt.Errorf("%v: cannot load events", err)
	}
	var user = onlyuser.User{}
	for _, event := range events {
		user.When(event)
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

func loadEvents(dbEvents []eventstore.DBEventItem) ([]interface{}, error) {
	var events []interface{}
	for _, dbEvent := range dbEvents {
		if dbEvent.Type == "UserCreated" {
			var e onlyuser.UserCreated
			err := json.Unmarshal([]byte(dbEvent.Data), &e)
			if err != nil {
				return nil, err
			}
			events = append(events, e)
		}
		if dbEvent.Type == "UserEmailChanged" {
			var e onlyuser.UserEmailChanged
			err := json.Unmarshal([]byte(dbEvent.Data), &e)
			if err != nil {
				return nil, err
			}
			events = append(events, e)
		}
	}
	return events, nil
}

func toKey(user onlyuser.User, i int) eventstore.EventKey {
	return eventstore.EventKey{
		ID:      fmt.Sprintf("user-%s", user.ID),
		Version: user.Version + i + 1,
	}
}
