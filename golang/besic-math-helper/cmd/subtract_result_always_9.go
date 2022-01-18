package cmd

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"math/rand"
	"time"
)

const (
	MIN    = 11
	MAX    = 20
	RESULT = 9
)

func SubtractResultAlways9(ctx *cli.Context) error {
	fmt.Println("Subtract result always 9!")
	rand.Seed(time.Now().UnixNano())
	f := rand.Intn(MAX-MIN+1) + MIN
	s := f - RESULT
	fmt.Printf("%d - %d =", f, s)
	return nil
}
