package main

import (
	"besic-math-helper/cmd"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func main() {
	app := &cli.App{
		Name: "creates basic math exercises",
		Commands: []*cli.Command{
			{
				Name:   "subtract_result_always_9",
				Usage:  "generates exercises with subtract operations that will always result in 9",
				Action: cmd.SubtractResultAlways9,
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatalf("error running the app: %v", err)
	}
}
