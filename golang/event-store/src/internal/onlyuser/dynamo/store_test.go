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

		// when
		actual := store.Load(ctx, user.ID)

		// then
		assert.Equal(t, user.ID, actual.ID)
		assert.Equal(t, user.Name, actual.Name)
		assert.Equal(t, user.Email, actual.Email)
	})

	t.Run("multiple events - success", func(t *testing.T) {
		db, table := dynamo.SetupTable(t, "EventStore")
		store := onlyuserdynamo.NewStore(db, table)

		// TODO: this shoudl be NewUser instead!
		user := onlyuser.User{
			ID:    "123456",
			Email: "sebastian@example.com",
		}
		user.ChangeEmail("seba@example.com")

		err := store.Save(ctx, user)
		assert.NoError(t, err)

		actual := store.Load(ctx, "123456")
		assert.Equal(t, "seba@example.com", actual.Email)
	})

	t.Run("test for 1MB return from Query", func(t *testing.T) {

	})
}
