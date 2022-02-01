package cmd

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"math/rand"
	"time"
)

func RandomAddition(ctx *cli.Context) error {
	fmt.Println("Random addition!")

	rand.Seed(time.Now().UnixNano())

	num := ctx.Int("num_exercises")
	exercises := generateRandomAdditionExercises(num)
	printExercises(exercises)
	return nil
}

func generateRandomAdditionExercises(num int) []string {
	exercises := make([]string, 0)
	for i := 0; i < num; i++ {
		f := rand.Intn(MAX + 1)
		s := rand.Intn(MAX - f + 1)
		exercises = append(exercises, fmt.Sprintf("%2d + %2d =", f, s))
	}
	return exercises
}
