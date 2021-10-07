package main

import (
	"github.com/urfave/cli/v2"
	"goclass/cmd"
	"log"
	"os"
)

func main() {
	app := cli.App{
		Name: "go class example service",
		Commands: []*cli.Command{
			{
				Name:   "hello",
				Usage:  "says hello",
				Action: cmd.RunHello,
			},
			{
				Name:   "server",
				Usage:  "starts the http server",
				Action: cmd.RunServer,
			},
			{
				Name:   "refresher",
				Usage:  "refreshes the local data",
				Action: cmd.RunRefresher,
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal("error running the app: %v", err)
	}
}
