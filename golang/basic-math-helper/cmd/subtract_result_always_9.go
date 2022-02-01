package cmd

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"math/rand"
	"time"
)

const (
	MIN    = 11
	MAX    = 30
	RESULT = 9
)

func SubtractResultAlways9(ctx *cli.Context) error {
	fmt.Println("Subtract result always 9!")

	rand.Seed(time.Now().UnixNano())

	num := ctx.Int("num_exercises")
	exercises := generateSubtractAlways9Exercise(num)
	printExercises(exercises)
	return nil
}

func generateSubtractAlways9Exercise(num int) []string {
	exercises := make([]string, 0)
	for i := 0; i < num; i++ {
		f := rand.Intn(MAX-MIN+1) + MIN
		s := f - RESULT
		exercises = append(exercises, fmt.Sprintf("%d - %2d =", f, s))
	}
	return exercises
}

func printExercises(exercises []string) {
	nextLine := 10
	for i, e := range exercises {
		fmt.Printf("%-14s", e)
		if i > 0 && i%nextLine == nextLine-1 {
			fmt.Println()
		}
	}
}
