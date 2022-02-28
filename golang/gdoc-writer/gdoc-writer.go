package main

import (
	"gdoc-writer/google"
	"gdoc-writer/stdin"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func main() {
	app := &cli.App{
		Name:        "gdoc-writer",
		Usage:       "creates Google Document with the content from stdin",
		Description: "creates Google Document with the content from stdin",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "prefix",
				Aliases: []string{"p"},
				Usage:   "prefix for the document title",
			},
		},
		Action: createDocumentFromStdin,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatalf("error running the app: %v", err)
	}
}

func createDocumentFromStdin(context *cli.Context) error {
	b, err := stdin.ReadStdin()
	if err != nil {
		log.Fatalf("Unable to read data from stdin: %v", err)
	}

	prefix := context.String("prefix")
	service := google.NewDocumentService()
	docLocation := google.CreateDocument(service, string(b), prefix)
	log.Printf("Successfully created Google Document under: %s", docLocation)
	return nil
}
