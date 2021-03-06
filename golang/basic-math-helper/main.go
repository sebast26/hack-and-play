package main

import (
	"basic-math-helper/cmd"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func main() {
	app := &cli.App{
		Name:        "basic-math-helper",
		Description: "generates basic math exercises",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:    "num_exercises",
				Aliases: []string{"n"},
				Value:   30,
				Usage:   "number of exercises to generate",
			},
		},
		Commands: []*cli.Command{
			{
				Name:   "subtract_result_always_9",
				Usage:  "generates exercises with subtract operations that will always result in 9",
				Action: cmd.SubtractResultAlways9,
			},
			{
				Name:   "random_subtract",
				Usage:  "generates random set of subtract exercises without any scheme",
				Action: cmd.RandomSubtract,
			},
			{
				Name:   "random_addition",
				Usage:  "generates random set of addition exercises without any scheme",
				Action: cmd.RandomAddition,
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatalf("error running the app: %v", err)
	}
}
