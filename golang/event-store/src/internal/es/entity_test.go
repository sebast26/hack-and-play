package es_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"sgorecki.me/golang/event-store/src/internal/es"
	eventstore "sgorecki.me/golang/event-store/src/internal/es/dynamo"
)

type EntityA struct {
	es.Entity
	content string
}

type EntityB struct {
	es.Entity
	content string
}

func (a *EntityA) When(event eventstore.DBEventItem) {
	a.content = "a when"
}

func (b *EntityB) When(event eventstore.DBEventItem) {
	b.content = "b when"
}

func TestEventSourcerEntity(t *testing.T) {
	t.Run("should invoke custom When() behaviour", func(t *testing.T) {
		// given
		a := EntityA{}
		b := EntityB{}

		// when
		a.When(eventstore.DBEventItem{})
		b.When(eventstore.DBEventItem{})

		// then
		assert.Equal(t, "a when", a.content)
		assert.Equal(t, "b when", b.content)
	})

	t.Run("should invoke custom When() through common method", func(t *testing.T) {
		// given
		a := EntityA{}
		b := EntityB{}

		// when
		a.Apply(&a, eventstore.DBEventItem{})
		b.Apply(&b, eventstore.DBEventItem{})

		// then
		assert.Equal(t, "a when", a.content)
		assert.Equal(t, "b when", b.content)
	})

	t.Run("should properly version entities - apply does not change entity version", func(t *testing.T) {
		// given
		a := EntityA{}

		// when
		a.Apply(&a, eventstore.DBEventItem{EventKey: eventstore.EventKey{Version: 0}})
		a.Apply(&a, eventstore.DBEventItem{EventKey: eventstore.EventKey{Version: 0}})
		a.Apply(&a, eventstore.DBEventItem{EventKey: eventstore.EventKey{Version: 0}})

		// then
		assert.Equal(t, 0, a.Version)
	})
}
