package cmd

import (
	"github.com/urfave/cli/v2"
)

func RunRefresher(ctx *cli.Context) error {
	s := getService()
	s.Refresh()

	return nil
}
