package complex

import "fmt"

type Point struct {
	X, Y int
}

type Circle struct {
	Center Point
	Radius int
}

type Wheel struct {
	Circle Circle
	Spokes int
}

func notSoGood() {
	var w Wheel
	w.Circle.Center.X = 8
	w.Circle.Center.Y = 8
	w.Circle.Radius = 5
	w.Spokes = 20
}

type BetterCircle struct {
	Point
	Radius int
}

type BetterWheel struct {
	BetterCircle
	Spokes int
}

func betterStructs() {
	var w BetterWheel
	w.X = 9
	w.Y = 9
	w.Radius = 5
	w.Spokes = 20

	// but still we need to...
	w = BetterWheel{
		BetterCircle: BetterCircle{
			Point: Point{
				X: 8,
				Y: 8,
			},
			Radius: 5,
		},
		Spokes: 20,
	}

	fmt.Println("%#v\n", w)
}
