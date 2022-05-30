package interfaces

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMultiTierSort(t *testing.T) {

	t.Run("single key", func(t *testing.T) {

		t.Run("sort by title", func(t *testing.T) {
			var byTitle = multiTierSort{t: freshTracks(), keys: []string{"title"}}

			sort.Sort(byTitle)

			assert.Equal(t, "Go", byTitle.t[0].Title)
			assert.Equal(t, "Go", byTitle.t[1].Title)
			assert.Equal(t, "Go Ahead", byTitle.t[2].Title)
			assert.Equal(t, "Ready 2 Go", byTitle.t[3].Title)
		})

		t.Run("sort by artist", func(t *testing.T) {
			var byArtist = multiTierSort{t: freshTracks(), keys: []string{"artist"}}

			sort.Sort(byArtist)

			assert.Equal(t, "Alicia Keys", byArtist.t[0].Artist)
			assert.Equal(t, "Delilah", byArtist.t[1].Artist)
			assert.Equal(t, "Martin Solveig", byArtist.t[2].Artist)
			assert.Equal(t, "Moby", byArtist.t[3].Artist)
		})

		t.Run("sort by album", func(t *testing.T) {
			var byAlbum = multiTierSort{t: freshTracks(), keys: []string{"album"}}

			sort.Sort(byAlbum)

			assert.Equal(t, "As I Am", byAlbum.t[0].Album)
			assert.Equal(t, "From the Roots Up", byAlbum.t[1].Album)
			assert.Equal(t, "Moby", byAlbum.t[2].Album)
			assert.Equal(t, "Smash", byAlbum.t[3].Album)
		})

		t.Run("sort by length", func(t *testing.T) {
			var byLength = multiTierSort{t: freshTracks(), keys: []string{"length"}}

			sort.Sort(byLength)

			assert.Equal(t, Length("3m37s"), byLength.t[0].Length)
			assert.Equal(t, Length("3m38s"), byLength.t[1].Length)
			assert.Equal(t, Length("4m24s"), byLength.t[2].Length)
			assert.Equal(t, Length("4m36s"), byLength.t[3].Length)
		})

		t.Run("sort by year", func(t *testing.T) {
			var byYear = multiTierSort{t: freshTracks(), keys: []string{"year"}}

			sort.Sort(byYear)

			assert.Equal(t, 1992, byYear.t[0].Year)
			assert.Equal(t, 2007, byYear.t[1].Year)
			assert.Equal(t, 2011, byYear.t[2].Year)
			assert.Equal(t, 2012, byYear.t[3].Year)
		})

	})

	t.Run("multiple keys", func(t *testing.T) {

		t.Run("sort by title and year", func(t *testing.T) {
			var byTitleYear = multiTierSort{t: freshTracks(), keys: []string{"title", "year"}}

			sort.Sort(byTitleYear)

			assert.Equal(t, "Go", byTitleYear.t[0].Title)
			assert.Equal(t, 1992, byTitleYear.t[0].Year)
			assert.Equal(t, "Go", byTitleYear.t[1].Title)
			assert.Equal(t, 2012, byTitleYear.t[1].Year)
			assert.Equal(t, "Go Ahead", byTitleYear.t[2].Title)
			assert.Equal(t, "Ready 2 Go", byTitleYear.t[3].Title)
		})

		t.Run("sort by title and album", func(t *testing.T) {
			var byTitleAlbum = multiTierSort{t: freshTracks(), keys: []string{"title", "album"}}

			sort.Sort(byTitleAlbum)

			assert.Equal(t, "Go", byTitleAlbum.t[0].Title)
			assert.Equal(t, "From the Roots Up", byTitleAlbum.t[0].Album)
			assert.Equal(t, "Go", byTitleAlbum.t[1].Title)
			assert.Equal(t, "Moby", byTitleAlbum.t[1].Album)
			assert.Equal(t, "Go Ahead", byTitleAlbum.t[2].Title)
			assert.Equal(t, "Ready 2 Go", byTitleAlbum.t[3].Title)
		})

	})

}

func freshTracks() []*Track {
	return []*Track{
		{"Go", "Delilah", "From the Roots Up", 2012, Length("3m38s")},
		{"Go", "Moby", "Moby", 1992, Length("3m37s")},
		{"Go Ahead", "Alicia Keys", "As I Am", 2007, Length("4m36s")},
		{"Ready 2 Go", "Martin Solveig", "Smash", 2011, Length("4m24s")},
	}
}
