package es_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"sgorecki.me/golang/event-store/src/internal/es"
)

type EntityA struct {
	es.Entity
	content string
}

type EntityB struct {
	es.Entity
	content string
}

func (a *EntityA) When(event interface{}) {
	a.content = "a when"
}

func (b *EntityB) When(event interface{}) {
	b.content = "b when"
}

func TestEventSourcerEntity(t *testing.T) {
	t.Run("should invoke custom When() behaviour", func(t *testing.T) {
		// given
		a := EntityA{}
		b := EntityB{}

		// when
		a.When("empty event")
		b.When("empty event")

		// then
		assert.Equal(t, "a when", a.content)
		assert.Equal(t, "b when", b.content)
	})

	t.Run("should invoke custom When() through common method", func(t *testing.T) {
		// given
		a := EntityA{}
		b := EntityB{}

		// when
		a.Apply(&a, "some event")
		b.Apply(&b, "some other event")

		// then
		assert.Equal(t, "a when", a.content)
		assert.Equal(t, "b when", b.content)
	})

	t.Run("should properly version entities - apply does not change entity version", func(t *testing.T) {
		// given
		a := EntityA{}

		// when
		a.Apply(&a, "version 0")
		a.Apply(&a, "version 0")
		a.Apply(&a, "version 0")

		// then
		assert.Equal(t, 0, a.Version)
	})
}
