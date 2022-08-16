package es

type shape struct {
	color string
}

type Circle struct {
	shape
	Radius int
}

func testInh() {
	circle := Circle{
		shape: shape{
			color: "black",
		},
		Radius: 10,
	}
	_ = circle.color
}
