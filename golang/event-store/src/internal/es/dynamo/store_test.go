package dynamo_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"git.naspersclassifieds.com/olxeu/specialized/kuna/platform-v2/testing/dynamo"
	eventstore "sgorecki.me/golang/event-store/src/internal/es/dynamo"
)

func TestEventStore(t *testing.T) {
	ctx := context.Background()

	t.Run("failure - concurrent update", func(t *testing.T) {
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
}