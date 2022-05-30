package interfaces

type multiTierSort struct {
	t    []*Track
	keys []string
}

func (x multiTierSort) Len() int      { return len(x.t) }
func (x multiTierSort) Swap(i, j int) { x.t[i], x.t[j] = x.t[j], x.t[i] }
func (x multiTierSort) Less(i, j int) bool {
	for _, k := range x.keys {
		if k == "title" {
			if x.t[i].Title != x.t[j].Title {
				return x.t[i].Title < x.t[j].Title
			}
		}
		if k == "artist" {
			if x.t[i].Artist != x.t[j].Artist {
				return x.t[i].Artist < x.t[j].Artist
			}
		}
		if k == "album" {
			if x.t[i].Album != x.t[j].Album {
				return x.t[i].Album < x.t[j].Album
			}
		}
		if k == "year" {
			if x.t[i].Year != x.t[j].Year {
				return x.t[i].Year < x.t[j].Year
			}
		}
		if k == "length" {
			if x.t[i].Length != x.t[j].Length {
				return x.t[i].Length < x.t[j].Length
			}
		}
	}
	return false
}
