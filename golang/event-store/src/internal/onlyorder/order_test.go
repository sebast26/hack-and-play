package onlyorder_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"sgorecki.me/golang/event-store/src/internal/onlyorder"
)

func TestNewOrder(t *testing.T) {
	t.Run("should create new Order", func(t *testing.T) {
		// when
		order := onlyorder.NewOrder()

		// then
		assert.NotEmpty(t, order.ID)
		assert.Empty(t, order.OrderItems)
		assert.Equal(t, 0, order.TotalAmount)
	})

	t.Run("should add item to Order", func(t *testing.T) {
		// given
		order := onlyorder.NewOrder()
		item := onlyorder.OrderItem{Name: "some item", Total: 10}

		// when
		order.AddItem(item)

		// then
		assert.Contains(t, order.OrderItems, item)
		assert.Equal(t, 10, order.TotalAmount)
	})
}
