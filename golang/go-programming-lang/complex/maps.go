package complex

import (
	"fmt"
	"sort"
	"strconv"
)

func mapOps() {
	// create
	ages := make(map[string]int)
	ages = map[string]int{
		"alice":   31,
		"charlie": 34,
	}

	// get
	fmt.Println(ages["alice"])
	fmt.Println(ages["seba"]) // safe, it will just return 0 value

	// delete
	delete(ages, "alice")
}

func mapSorted() {
	ages := map[string]int{
		"charlie": 34,
		"alice":   31,
		"bob":     33,
	}

	// getting slice with names
	names := make([]string, 0, len(ages))
	for name := range ages {
		names = append(names, name)
	}
	// sorting slice (there is no sort map func)
	sort.Strings(names)
	// output sorted map
	for _, name := range names {
		fmt.Printf("%s\t%d\n", name, ages[name])
	}
}

func mapNoValue() {
	ages := map[string]int{
		"alice":   31,
		"charlie": 34,
	}

	// get
	fmt.Println(ages["seba"]) // safe, it will just return 0 value

	if age, present := ages["seba"]; !present {
		fmt.Println("No value for 'seba', age is " + strconv.Itoa(age))
	}

}

// instead of putting slice as key, we put string that represents this list
func k(list []string) string { return fmt.Sprintf("%q", list) }

func mapWithSliceAsKey(list []string) int {
	var m = make(map[string]int)

	// add
	m[k(list)]++

	// count
	return m[k(list)]
}

func RunMaps() {
	mapOps()
	mapSorted()
	mapNoValue()
}
