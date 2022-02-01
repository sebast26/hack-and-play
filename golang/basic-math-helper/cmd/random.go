package cmd

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"math/rand"
	"time"
)

func RandomSubtract(ctx *cli.Context) error {
	fmt.Println("Random subtract!")

	rand.Seed(time.Now().UnixNano())

	num := ctx.Int("num_exercises")
	exercises := generateRandomExercise(num)
	printExercises(exercises)
	return nil
}

func generateRandomExercise(num int) []string {
	exercises := make([]string, 0)
	for i := 0; i < num; i++ {
		f := rand.Intn(MAX + 1)
		s := rand.Intn(f + 1)
		exercises = append(exercises, fmt.Sprintf("%d - %2d =", f, s))
	}
	return exercises
}
