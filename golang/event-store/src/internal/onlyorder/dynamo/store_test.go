package dynamo_test

import (
	"context"
	"testing"

	"git.naspersclassifieds.com/olxeu/specialized/kuna/platform-v2/testing/dynamo"
	"github.com/stretchr/testify/assert"
	eventstore "sgorecki.me/golang/event-store/src/internal/es/dynamo"
	"sgorecki.me/golang/event-store/src/internal/onlyorder"
	orderstore "sgorecki.me/golang/event-store/src/internal/onlyorder/dynamo"
)

// Only basic tests for order (not checking versioning, concurrency, limits, etc). Original tests in onlyuser.
func TestEventStore(t *testing.T) {
	ctx := context.Background()

	t.Run("success - create order", func(t *testing.T) {
		// given
		db, table := dynamo.SetupTable(t, "EventStore")
		store := orderstore.NewStore(eventstore.NewStore(db, table))
		order := onlyorder.NewOrder()

		// when
		err := store.Save(ctx, order)

		// then
		assert.NoError(t, err)

		// when
		actual, err := store.Load(ctx, order.ID)

		// then
		assert.NoError(t, err)
		assert.Equal(t, order.ID, actual.ID)
		assert.Empty(t, actual.OrderItems)
		assert.Equal(t, 0, actual.TotalAmount)
		assert.Equal(t, 1, actual.Version) // internals: should it be here?
	})

	t.Run("success - adding item to order", func(t *testing.T) {
		// given
		db, table := dynamo.SetupTable(t, "EventStore")
		store := orderstore.NewStore(eventstore.NewStore(db, table))
		order := onlyorder.NewOrder()
		item := onlyorder.OrderItem{Name: "shiny item", Total: 1000}

		// when
		order.AddItem(item)
		err := store.Save(ctx, order)

		// then
		assert.NoError(t, err)
		actual, err := store.Load(ctx, order.ID)
		assert.NoError(t, err)
		assert.Contains(t, actual.OrderItems, item)
		assert.Equal(t, 1000, actual.TotalAmount)
		assert.Equal(t, 2, actual.Version) // internals: should it be here?
	})
}
