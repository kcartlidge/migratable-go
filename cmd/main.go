// Main entry point.
// AGPL license. Copyright 2024 K Cartlidge.

package main

import (
	"fmt"
)

func main() {
	// Welcome.
	fmt.Println()
	fmt.Println("MIGRATABLE")
	fmt.Println()

	// Gather the user request.
	c, err := getConfig()
	if err == nil {
		// Action it.
		c.describe()
	}

	// Handle any errors bubbled up.
	if err != nil {
		fmt.Println()
		fmt.Println("ERROR")
		fmt.Println(err.Error())
		fmt.Println()
	} else {
		fmt.Println()
		fmt.Println("Done.")
		fmt.Println()
	}
}
