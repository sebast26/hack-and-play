package es

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type BaseEntity interface {
	Apply() string
	When() string
}

type abstractBaseEntity struct {
	BaseEntity
}

func (e abstractBaseEntity) Apply() string {
	return e.When()
}

type User struct {
	abstractBaseEntity
}

func NewUser() User {
	user := User{}
	user.abstractBaseEntity.BaseEntity = user
	return user
}

func (u User) When() string {
	return "user"
}

type Order struct {
	abstractBaseEntity
}

func NewOrder() Order {
	order := Order{}
	order.abstractBaseEntity.BaseEntity = order
	return order
}

func (o Order) When() string {
	return "order"
}

type Customer interface {
	When() string
}

type ESEntity struct {
	ID string
}

type RushCustomer struct {
	ESEntity
}

type SlowCustomer struct {
	ESEntity
}

func (c *RushCustomer) When() string {
	return "rush"
}

func (c *SlowCustomer) When() string {
	return "slow"
}

func (e *ESEntity) When() string {
	return "es"
}

func (e *ESEntity) Apply() string {
	return e.When()
}

type iEntity interface {
	Apply(iEntity) string
	When() string
}

type baseEntity struct {
	name string
}

func (e *baseEntity) Apply(entity iEntity) string {
	return entity.When()
}

type Booking struct {
	baseEntity
}

type Calendar struct {
	baseEntity
}

func (b *Booking) When() string {
	fmt.Println(b.name)
	return "booking"
}

func (c *Calendar) When() string {
	fmt.Println(c.name)
	return "calendar"
}

func TestInheritanceInGo(t *testing.T) {

	t.Run("test custom behaviour", func(t *testing.T) {
		rushCustomer := RushCustomer{}
		slowCustomer := SlowCustomer{}

		assert.Equal(t, "rush", rushCustomer.When())
		assert.Equal(t, "slow", slowCustomer.When())
	})

	t.Run("test custom behaviour for common method", func(t *testing.T) {
		t.Skip("it will fail with this implementation")
		rushCustomer := RushCustomer{}
		slowCustomer := SlowCustomer{}

		assert.Equal(t, "rush", rushCustomer.Apply())
		assert.Equal(t, "slow", slowCustomer.Apply())
	})

	t.Run("2nd impl - test custom behaviour", func(t *testing.T) {
		user := User{}
		order := Order{}

		assert.Equal(t, "user", user.When())
		assert.Equal(t, "order", order.When())
	})

	t.Run("2nd impl - test custom behaviour for common method", func(t *testing.T) {
		user := NewUser()
		order := NewOrder()

		// in this approach it is not possible to use *User receiver and modify struct in place
		assert.Equal(t, "user", user.Apply())
		assert.Equal(t, "order", order.Apply())
	})

	t.Run("3rd impl - test custom behaviour", func(t *testing.T) {
		booking := Booking{}
		calendar := Calendar{}

		assert.Equal(t, "booking", booking.When())
		assert.Equal(t, "calendar", calendar.When())
	})

	t.Run("3rd impl - test custom behaviour for common method", func(t *testing.T) {
		booking := Booking{baseEntity{name: "booking"}}
		calendar := Calendar{baseEntity{name: "calendar"}}

		assert.Equal(t, "booking", booking.Apply(&booking))
		assert.Equal(t, "calendar", calendar.Apply(&calendar))
	})
}
