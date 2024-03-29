package dynamo_test

import (
	"context"
	"testing"
	"time"

	"git.naspersclassifieds.com/olxeu/specialized/kuna/platform-v2/testing/dynamo"
	"github.com/stretchr/testify/assert"
	eventstore "sgorecki.me/golang/event-store/src/internal/es/dynamo"
	"sgorecki.me/golang/event-store/src/internal/onlyuser"
	userstore "sgorecki.me/golang/event-store/src/internal/onlyuser/dynamo"
)

func TestEventStore(t *testing.T) {
	ctx := context.Background()
	realClock := func() time.Time { return time.Now() }

	t.Run("success - create user", func(t *testing.T) {
		// given
		db, table := dynamo.SetupTable(t, "EventStore")
		store := userstore.NewStore(eventstore.NewStore(db, table, realClock))
		user := onlyuser.NewUser("name", "email")

		// when
		err := store.Save(ctx, user)

		// then
		assert.NoError(t, err)

		// when
		actual, err := store.Load(ctx, user.ID)

		// then
		assert.NoError(t, err)
		assert.Equal(t, user.ID, actual.ID)
		assert.Equal(t, user.Name, actual.Name)
		assert.Equal(t, user.Email, actual.Email)
		assert.Equal(t, 1, actual.Version) // internals: should it be here?
	})

	t.Run("success - user changed email", func(t *testing.T) {
		// given
		db, table := dynamo.SetupTable(t, "EventStore")
		store := userstore.NewStore(eventstore.NewStore(db, table, realClock))
		user := onlyuser.NewUser("name", "email")

		// when
		user.ChangeEmail("otheremail")
		err := store.Save(ctx, user)

		// then
		assert.NoError(t, err)
		actual, err := store.Load(ctx, user.ID)
		assert.NoError(t, err)
		assert.Equal(t, "name", actual.Name)
		assert.Equal(t, "otheremail", actual.Email)
		assert.Equal(t, 2, actual.Version) // internals: should it be here?
	})

	t.Run("success - user changed email, order of events matter", func(t *testing.T) {
		// given
		db, table := dynamo.SetupTable(t, "EventStore")
		store := userstore.NewStore(eventstore.NewStore(db, table, realClock))
		user := onlyuser.NewUser("name", "email")

		// when
		user.ChangeEmail("otheremail")
		user.ChangeEmail("finalemail")
		err := store.Save(ctx, user)

		// then
		assert.NoError(t, err)
		actual, err := store.Load(ctx, user.ID)
		assert.NoError(t, err)
		assert.Equal(t, "name", actual.Name)
		assert.Equal(t, "finalemail", actual.Email)
		assert.Equal(t, 3, actual.Version)
	})

	t.Run("success - multiple users and mutations", func(t *testing.T) {
		// given
		// TODO!!

		// when

		// then

	})

	t.Run("test for 1MB return from Query", func(t *testing.T) {

	})

	// TODO: duplicated tests from es/dynamo
	t.Run("failure - concurrent update", func(t *testing.T) {
		// given
		db, table := dynamo.SetupTable(t, "EventStore")
		store := userstore.NewStore(eventstore.NewStore(db, table, realClock))
		user := onlyuser.NewUser("name", "email")

		// when
		err1 := store.Save(ctx, user)
		err2 := store.Save(ctx, user)

		// then
		assert.NoError(t, err1)
		assert.Error(t, err2)
		assert.ErrorIs(t, err2, eventstore.ErrConcurrentUpdate)
	})
}
