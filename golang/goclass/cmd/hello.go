package cmd

import (
	"fmt"
	"github.com/urfave/cli/v2"
)

func RunHello(ctx *cli.Context) error {
	fmt.Println("Hello!")
	return nil
}
