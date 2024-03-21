// String routines.
// AGPL license. Copyright 2024 K Cartlidge.

package main

import (
	"fmt"
	"strings"
)

// section writes a section header to the output.
func section(content string) {
	minwid := 40
	wid := len(content)
	if wid < minwid {
		wid = minwid
	}
	fmt.Println()
	fmt.Println()
	fmt.Println(content)
	fmt.Println(strings.Repeat("=", wid))
}

// limitTo forces the string to fit the length constraint provided.
func limitTo(original string, max int) string {
	i := 0
	for j := range original {
		if i == max {
			return original[:j]
		}
		i++
	}
	return original
}
