package es_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"sgorecki.me/golang/event-store/src/internal/es"
)

type EntityA struct {
	es.Entity
}

type EntityB struct {
	es.Entity
}

func (a *EntityA) When(event interface{}) {
	a.Version = 1
}

func (b *EntityB) When(event interface{}) {
	b.Version = 2
}

func TestEventSourcerEntity(t *testing.T) {
	t.Run("should invoke custom When() behaviour", func(t *testing.T) {
		// given
		a := EntityA{Entity: es.Entity{Version: 0}}
		b := EntityB{Entity: es.Entity{Version: 0}}

		// when
		a.When("empty event")
		b.When("empty event")

		// then
		assert.Equal(t, 1, a.Version)
		assert.Equal(t, 2, b.Version)
	})

	t.Run("should invoke custom When() through common method", func(t *testing.T) {
		// given
		a := EntityA{Entity: es.Entity{Version: 0}}
		b := EntityB{Entity: es.Entity{Version: 0}}

		// when
		a.Apply(&a, "some event")
		b.Apply(&b, "some other event")

		// then
		assert.Equal(t, 1, a.Version)
		assert.Equal(t, 2, b.Version)
	})
}
