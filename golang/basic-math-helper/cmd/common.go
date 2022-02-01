package cmd

import "fmt"

func printExercises(exercises []string) {
	nextLine := 8
	for i, e := range exercises {
		fmt.Printf("%-14s", e)
		if i > 0 && i%nextLine == nextLine-1 {
			fmt.Println()
		}
	}
	fmt.Println()
}
