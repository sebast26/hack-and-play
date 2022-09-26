package dynamo_test

import (
	"context"
	"strings"
	"testing"

	"git.naspersclassifieds.com/olxeu/specialized/kuna/platform-v2/testing/dynamo"
	"github.com/stretchr/testify/assert"
	eventstore "sgorecki.me/golang/event-store/src/internal/es/dynamo"
)

func TestEventStore(t *testing.T) {
	ctx := context.Background()

	t.Run("success - store entity with single event", func(t *testing.T) {
		// given
		db, table := dynamo.SetupTable(t, "EventStore")
		store := eventstore.NewStore(db, table)
		item := eventstore.DBEventItem{
			EventKey: eventstore.EventKey{
				ID:      "user-1",
				Version: 1,
			},
		}

		// when
		err := store.AppendEvents(ctx, []eventstore.DBEventItem{item})

		// then
		assert.NoError(t, err)

		// when
		actual, err := store.ReadEvents(ctx, item.ID)
		assert.NoError(t, err)
		assert.Len(t, actual, 1)
		assert.Equal(t, "user-1", actual[0].ID)
		assert.Equal(t, 1, actual[0].Version)
	})

	t.Run("success - no events", func(t *testing.T) {
		// given
		db, table := dynamo.SetupTable(t, "EventStore")
		store := eventstore.NewStore(db, table)

		// when
		err := store.AppendEvents(ctx, []eventstore.DBEventItem{})

		// then
		assert.NoError(t, err)
	})

	t.Run("success - multiple events", func(t *testing.T) {
		// given
		db, table := dynamo.SetupTable(t, "EventStore")
		store := eventstore.NewStore(db, table)
		items := []eventstore.DBEventItem{
			{
				EventKey: eventstore.EventKey{ID: "stream-1", Version: 1},
			},
			{
				EventKey: eventstore.EventKey{ID: "stream-1", Version: 2},
			},
		}

		// when
		err := store.AppendEvents(ctx, items)

		// then
		assert.NoError(t, err)

		// and
		actual, err := store.ReadEvents(ctx, "stream-1")
		assert.NoError(t, err)
		assert.Len(t, actual, 2)
		assert.Equal(t, eventstore.DBEventItem{EventKey: eventstore.EventKey{ID: "stream-1", Version: 1}}, actual[0])
		assert.Equal(t, eventstore.DBEventItem{EventKey: eventstore.EventKey{ID: "stream-1", Version: 2}}, actual[1])
	})

	t.Run("failure - invalid event in batch, no entries added", func(t *testing.T) {
		// given
		db, table := dynamo.SetupTable(t, "EventStore")
		store := eventstore.NewStore(db, table)
		item := eventstore.DBEventItem{EventKey: eventstore.EventKey{ID: "stream-1", Version: 1}}
		invalidItem := eventstore.DBEventItem{EventKey: eventstore.EventKey{}}

		// when
		err := store.AppendEvents(ctx, []eventstore.DBEventItem{item, invalidItem})

		// then
		assert.Error(t, err)

		// and
		actual, err := store.ReadEvents(ctx, "stream-1")
		assert.NoError(t, err)
		assert.Len(t, actual, 0)
	})

	t.Run("concurrent update", func(t *testing.T) {

		t.Run("failure - adding same event", func(t *testing.T) {
			// given
			db, table := dynamo.SetupTable(t, "EventStore")
			store := eventstore.NewStore(db, table)
			item := eventstore.DBEventItem{
				EventKey: eventstore.EventKey{
					ID:      "user-1",
					Version: 1,
				},
			}

			// when
			err := store.AppendEvents(ctx, []eventstore.DBEventItem{item})
			err2 := store.AppendEvents(ctx, []eventstore.DBEventItem{item})

			// then
			assert.NoError(t, err)
			assert.Error(t, err2)
			assert.ErrorIs(t, err2, eventstore.ErrConcurrentUpdate)
		})

		t.Run("failure - concurrent update in second batch", func(t *testing.T) {
			// given
			db, table := dynamo.SetupTable(t, "EventStore")
			store := eventstore.NewStore(db, table)
			itemV1 := eventstore.DBEventItem{EventKey: eventstore.EventKey{ID: "stream-1", Version: 1}}
			itemV2 := eventstore.DBEventItem{EventKey: eventstore.EventKey{ID: "stream-1", Version: 2}}

			// when
			err := store.AppendEvents(ctx, []eventstore.DBEventItem{itemV1})

			// then
			assert.NoError(t, err)

			// when
			err = store.AppendEvents(ctx, []eventstore.DBEventItem{itemV2, itemV1})

			// then
			assert.Error(t, err)

			// when
			actual, err := store.ReadEvents(ctx, "stream-1")
			assert.Len(t, actual, 1)
			assert.Equal(t, itemV1, actual[0])
		})

		t.Run("success - unique partition key + sort key", func(t *testing.T) {
			// given
			db, table := dynamo.SetupTable(t, "EventStore")
			store := eventstore.NewStore(db, table)
			itemV1 := eventstore.DBEventItem{EventKey: eventstore.EventKey{ID: "stream-1", Version: 1}}
			otherV1 := eventstore.DBEventItem{EventKey: eventstore.EventKey{ID: "stream-2", Version: 1}}

			// when
			itemErr := store.AppendEvents(ctx, []eventstore.DBEventItem{itemV1})
			otherErr := store.AppendEvents(ctx, []eventstore.DBEventItem{otherV1})

			// then
			assert.NoError(t, itemErr)
			assert.NoError(t, otherErr)
		})
	})

	t.Run("item collection limit of 1MB", func(t *testing.T) {
		t.Run("success - paging through events with over 1MB items", func(t *testing.T) {
			// given
			db, table := dynamo.SetupTable(t, "EventStore")
			store := eventstore.NewStore(db, table)
			items := []eventstore.DBEventItem{
				generateItem(1, generate100KBString()),
				generateItem(2, generate100KBString()),
				generateItem(3, generate100KBString()),
				generateItem(4, generate100KBString()),
				generateItem(5, generate100KBString()),
				generateItem(6, generate100KBString()),
				generateItem(7, generate100KBString()),
				generateItem(8, generate100KBString()),
				generateItem(9, generate100KBString()),
				generateItem(10, generate100KBString()),
				generateItem(11, generate100KBString()),
				generateItem(12, generate100KBString()),
			}

			// when
			err := store.AppendEvents(ctx, items)

			// then
			assert.NoError(t, err)

			// and
			actual, err := store.ReadEvents(ctx, "event-1")
			assert.NoError(t, err)
			assert.Len(t, actual, 12)
		})
	})

}

func generate1KBString() string {
	var sb strings.Builder
	for i := 0; i < 1024; i++ {
		sb.WriteRune('a')
	}
	return sb.String()
}

func generate100KBString() string {
	var sb strings.Builder
	for i := 0; i < 100; i++ {
		sb.WriteString(generate1KBString())
	}
	return sb.String()
}

func generateItem(version int, data string) eventstore.DBEventItem {
	return eventstore.DBEventItem{
		EventKey: eventstore.EventKey{
			ID:      "event-1",
			Version: version,
		},
		Type: "test-event",
		Data: data,
	}
}
