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

	t.Run("save & get - success", func(t *testing.T) {
		db, table := dynamo.SetupTable(t, "EventStore")
		store := onlyuserdynamo.NewStore(db, table)

		user := onlyuser.User{
			ID:    "123456",
			Name:  "Sebastian",
			Email: "sebastian@example.com",
		}
		user.ChangeEmail("seba@example.com")

		err := store.Save(ctx, user)
		assert.NoError(t, err)

		actual := store.Load(ctx, "123456")
		assert.Equal(t, user.ID, actual.ID)
		assert.Equal(t, user.Name, actual.Name)
		assert.Equal(t, "seba@example.com", actual.Email)
	})
}
