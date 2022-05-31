package interfaces

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_sortedTracks_Sorting(t *testing.T) {

	t.Run("when empty sorting", func(t *testing.T) {
		st := sortedTracks{sorting: []string{}}

		assert.Equal(t, "?", st.Sorting())
	})

	t.Run("when single column", func(t *testing.T) {
		st := sortedTracks{sorting: []string{"title"}}

		assert.Equal(t, "?sort=title&", st.Sorting())
	})

	t.Run("when multiple columns", func(t *testing.T) {
		st := sortedTracks{sorting: []string{"title", "album", "year"}}

		assert.Equal(t, "?sort=title&sort=album&sort=year&", st.Sorting())
	})

}
