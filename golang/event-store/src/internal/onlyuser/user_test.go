package onlyuser_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"sgorecki.me/golang/event-store/src/internal/onlyuser"
)

func TestUser(t *testing.T) {
	t.Run("should create a user", func(t *testing.T) {
		user := onlyuser.NewUser("name", "email")

		assert.NotEmpty(t, user.ID)
		assert.Equal(t, user.Name, "name")
		assert.Equal(t, user.Email, "email")
	})

	t.Run("should update user's email", func(t *testing.T) {
		user := onlyuser.NewUser("name", "email")
		user.ChangeEmail("otheremail")

		assert.Equal(t, user.Email, "otheremail")
	})
}
