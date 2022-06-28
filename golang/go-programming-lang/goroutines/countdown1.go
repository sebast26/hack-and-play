package goroutines

import (
	"fmt"
	"time"
)

func Countdown1() {
	fmt.Println("Commencing countdown.")
	tick := time.Tick(1 * time.Second)
	for countdown := 10; countdown > 0; countdown-- {
		fmt.Println(countdown)
		<-tick
	}
	Launch()
}

func Launch() {
	fmt.Println("Launch!")
}
