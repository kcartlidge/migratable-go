// String routines.
// AGPL license. Copyright 2024 K Cartlidge.

package main

import (
	"strings"
)

const valid = "abcdefghijklmnopqrstuvwxyz+!&,()0123456789"

// slugify converts text to a URL-friendly slug.
func slugify(original string) string {
	r := ""
	l := ""
	for _, ru := range strings.ToLower(original) {
		c := string(ru)
		if strings.Index(valid, c) > -1 {
			r = r + c
			l = c
		} else if l != "-" {
			r = r + "-"
			l = "-"
		}
	}
	return strings.TrimPrefix(strings.TrimSuffix(r, "-"), "-")
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
