package interfaces

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"text/tabwriter"
	"time"
)

type Track struct {
	Title  string
	Artist string
	Album  string
	Year   int
	Length time.Duration
}

var tracks = []*Track{
	{"Go", "Delilah", "From the Roots Up", 2012, Length("3m38s")},
	{"Go", "Moby", "Moby", 1992, Length("3m37s")},
	{"Go Ahead", "Alicia Keys", "As I Am", 2007, Length("4m36s")},
	{"Ready 2 Go", "Martin Solveig", "Smash", 2011, Length("4m24s")},
}

func Length(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		panic(s)
	}
	return d
}

func printTracks(tracks []*Track) {
	const format = "%v\t%v\t%v\t%v\t%v\t\n"
	tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
	fmt.Fprintf(tw, format, "Title", "Artist", "Album", "Year", "Length")
	fmt.Fprintf(tw, format, "-----", "------", "-----", "----", "------")
	for _, t := range tracks {
		fmt.Fprintf(tw, format, t.Title, t.Artist, t.Album, t.Year, t.Length)
	}
	tw.Flush()
}

type byArtist []*Track

func (x byArtist) Len() int           { return len(x) }
func (x byArtist) Less(i, j int) bool { return x[i].Artist < x[j].Artist }
func (x byArtist) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

type customSort struct {
	t    []*Track
	less func(x, y *Track) bool
}

func (x customSort) Len() int           { return len(x.t) }
func (x customSort) Less(i, j int) bool { return x.less(x.t[i], x.t[j]) }
func (x customSort) Swap(i, j int)      { x.t[i], x.t[j] = x.t[j], x.t[i] }

func RunSorting() {
	fmt.Println("Not sorted tracks")
	printTracks(tracks)
	fmt.Println("Sorted by Artist")
	sort.Sort(byArtist(tracks))
	printTracks(tracks)
	fmt.Println("Reversed by Artist")
	sort.Sort(sort.Reverse(byArtist(tracks)))
	printTracks(tracks)
	fmt.Println("Custom sort: byTitle then byYear then byLength")
	sort.Sort(customSort{tracks, func(x, y *Track) bool {
		if x.Title != y.Title {
			return x.Title < y.Title
		}
		if x.Year != y.Year {
			return x.Year < y.Year
		}
		if x.Length != y.Length {
			return x.Length < y.Length
		}
		return false
	}})
	printTracks(tracks)
}

type sortedTracks struct {
	Tracks  []*Track
	sorting []string
}

func printTracksTable(w io.Writer, st sortedTracks) {
	tmpl := template.Must(template.ParseFiles("interfaces/tracks_table.html"))
	err := tmpl.Execute(w, st)
	if err != nil {
		log.Fatalf("%v: error preparing HTML table", err)
	}
}

func RunMultiTierTableSort() {
	http.HandleFunc("/tracks", multiTierSortTracks)
	http.ListenAndServe(":8099", nil)
}

func multiTierSortTracks(writer http.ResponseWriter, request *http.Request) {
	s := request.URL.Query()["sort"]
	st := sortedTracks{
		Tracks:  tracks,
		sorting: s,
	}

	multiSort := multiTierSort{t: st.Tracks, keys: st.sorting}
	sort.Sort(multiSort)

	printTracksTable(writer, st)
}

func (st sortedTracks) Sorting() string {
	var sb strings.Builder
	sb.WriteString("?")
	for _, s := range st.sorting {
		sb.WriteString("sort=" + s + "&")
	}
	return sb.String()
}
