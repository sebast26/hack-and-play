package dynamo_test

import (
	"context"
	"testing"
	"time"

	"git.naspersclassifieds.com/olxeu/specialized/kuna/platform-v2/testing/dynamo"
	"github.com/stretchr/testify/assert"
	eventstore "sgorecki.me/golang/event-store/src/internal/es/dynamo"
)

func TestEventStore(t *testing.T) {
	ctx := context.Background()
	fixedCreatedAt := "2022-09-27T10:15:30.000051234Z"
	fixedClock := func() time.Time { return time.Date(2022, 9, 27, 10, 15, 30, 51234, time.UTC) }

	t.Run("success - store entity with single event", func(t *testing.T) {
		// given
		db, table := dynamo.SetupTable(t, "EventStore")
		store := eventstore.NewStore(db, table, fixedClock)
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
		store := eventstore.NewStore(db, table, fixedClock)

		// when
		err := store.AppendEvents(ctx, []eventstore.DBEventItem{})

		// then
		assert.NoError(t, err)
	})

	t.Run("success - multiple events", func(t *testing.T) {
		// given
		db, table := dynamo.SetupTable(t, "EventStore")
		store := eventstore.NewStore(db, table, fixedClock)
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
		assert.Equal(t, eventstore.DBEventItem{EventKey: eventstore.EventKey{ID: "stream-1", Version: 1}, CreatedAt: fixedCreatedAt}, actual[0])
		assert.Equal(t, eventstore.DBEventItem{EventKey: eventstore.EventKey{ID: "stream-1", Version: 2}, CreatedAt: fixedCreatedAt}, actual[1])
	})

	t.Run("failure - invalid event in batch, no entries added", func(t *testing.T) {
		// given
		db, table := dynamo.SetupTable(t, "EventStore")
		store := eventstore.NewStore(db, table, fixedClock)
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

	t.Run("datetime of the events", func(t *testing.T) {
		t.Run("success - register time of the event", func(t *testing.T) {
			// given
			db, table := dynamo.SetupTable(t, "EventStore")
			store := eventstore.NewStore(db, table, fixedClock)
			item := eventstore.DBEventItem{EventKey: eventstore.EventKey{ID: "stream-1", Version: 1}}

			// when
			err := store.AppendEvents(ctx, []eventstore.DBEventItem{item})

			// then
			assert.NoError(t, err)
			actual, err := store.ReadEvents(ctx, "stream-1")
			assert.Equal(t, "2022-09-27T10:15:30.000051234Z", actual[0].CreatedAt)
		})

		t.Run("success - zero nano-seconds format", func(t *testing.T) {
			// given
			db, table := dynamo.SetupTable(t, "EventStore")
			zeroNanoClock := func() time.Time { return time.Date(2022, 9, 27, 15, 10, 29, 0, time.UTC) }
			expectedCreatedAt := "2022-09-27T15:10:29.000000000Z"
			store := eventstore.NewStore(db, table, zeroNanoClock)
			item := eventstore.DBEventItem{EventKey: eventstore.EventKey{ID: "stream-1", Version: 1}}

			// when
			err := store.AppendEvents(ctx, []eventstore.DBEventItem{item})

			// then
			assert.NoError(t, err)
			actual, err := store.ReadEvents(ctx, "stream-1")
			assert.Equal(t, expectedCreatedAt, actual[0].CreatedAt)
		})
	})

	t.Run("concurrent update", func(t *testing.T) {

		t.Run("failure - adding same event", func(t *testing.T) {
			// given
			db, table := dynamo.SetupTable(t, "EventStore")
			store := eventstore.NewStore(db, table, fixedClock)
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
			store := eventstore.NewStore(db, table, fixedClock)
			itemV1 := eventstore.DBEventItem{EventKey: eventstore.EventKey{ID: "stream-1", Version: 1}, CreatedAt: fixedCreatedAt}
			itemV2 := eventstore.DBEventItem{EventKey: eventstore.EventKey{ID: "stream-1", Version: 2}, CreatedAt: fixedCreatedAt}

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
			store := eventstore.NewStore(db, table, fixedClock)
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
			store := eventstore.NewStore(db, table, fixedClock)
			kb380 := string(make([]byte, 380*1024))
			items := []eventstore.DBEventItem{
				generateItem(1, kb380),
				generateItem(2, kb380),
				generateItem(3, kb380),
			}

			// when
			err := store.AppendEvents(ctx, items)

			// then
			assert.NoError(t, err)

			// and
			actual, err := store.ReadEvents(ctx, "event-1")
			assert.NoError(t, err)
			assert.Len(t, actual, 3)
		})
	})

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
