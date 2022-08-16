package dynamo_test

import (
	"context"
	"testing"

	"git.naspersclassifieds.com/olxeu/specialized/kuna/platform-v2/testing/dynamo"
	"github.com/stretchr/testify/assert"
	"sgorecki.me/golang/event-store/src/internal/onlyuser"
	onlyuserdynamo "sgorecki.me/golang/event-store/src/internal/onlyuser/dynamo"
)

func TestEventStore(t *testing.T) {
	ctx := context.Background()

	t.Run("success - create user", func(t *testing.T) {
		// given
		db, table := dynamo.SetupTable(t, "EventStore")
		store := onlyuserdynamo.NewStore(db, table)
		user := onlyuser.NewUser("name", "email")

		// when
		err := store.Save(ctx, user)

		// then
		assert.NoError(t, err)
		assert.Equal(t, 1, user.Version) // internals: should it be here?

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
		store := onlyuserdynamo.NewStore(db, table)
		user := onlyuser.NewUser("name", "email")

		// when
		user.ChangeEmail("otheremail")
		err := store.Save(ctx, user)

		// then
		assert.NoError(t, err)
		assert.Equal(t, 2, user.Version) // internals: should it be here?
		actual, err := store.Load(ctx, user.ID)
		assert.NoError(t, err)
		assert.Equal(t, "name", actual.Name)
		assert.Equal(t, "otheremail", actual.Email)
		assert.Equal(t, 2, actual.Version) // internals: should it be here?
	})

	t.Run("success - user changed email, order of events matter", func(t *testing.T) {
		// given
		db, table := dynamo.SetupTable(t, "EventStore")
		store := onlyuserdynamo.NewStore(db, table)
		user := onlyuser.NewUser("name", "email")

		// when
		user.ChangeEmail("otheremail")
		user.ChangeEmail("finalemail")
		err := store.Save(ctx, user)

		// then
		assert.NoError(t, err)
		assert.Equal(t, 3, user.Version)
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
}
