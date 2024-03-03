package main

import (
	"errors"
	"fmt"
	"path"
	"strconv"
	"unicode"
)

// migrations holds all loaded migrations.
type migrations map[int]migration

// migration holds one complete migration.
type migration struct {
	Name, Up, Down string
	Version        int
}

// loadMigrations reads in all the migrations in the
// sub-folders of the specified folder.
func loadMigrations(folder string) (*migrations, error) {
	ms := migrations{}

	fmt.Println("READING")

	// Scan for migrations folders.
	if ok, err := exists(folder); !ok || (err != nil) {
		return nil, errors.New("cannot read migrations folder")
	}

	folders, err := getSubFolders(folder)
	if err != nil {
		return nil, err
	}

	for _, f := range folders {
		fmt.Printf("  %v\n", f)

		// Extract the leading version number and the name.
		digits := []rune{}
		name := []rune{}
		inDigits := true
		for _, ch := range f {
			if inDigits && unicode.IsDigit(ch) {
				digits = append(digits, ch)
			} else {
				inDigits = false
				name = append(name, ch)
			}
		}
		if len(digits) == 0 {
			return nil, errors.New("migration folder missing version: " + f)
		}
		v, _ := strconv.Atoi(string(digits))

		// Populate the migration.
		subFolder := path.Join(folder, f)
		sqlUp, err := loadFile(subFolder, "up.sql", true)
		if err != nil {
			return nil, err
		}
		sqlDn, err := loadFile(subFolder, "down.sql", true)
		if err != nil {
			return nil, err
		}
		m := migration{
			Version: v,
			Name:    string(name),
			Up:      sqlUp,
			Down:    sqlDn,
		}

		// Gather by version.
		ms[v] = m
	}

	// Check the version numbers are complete and go from 1 upward.
	if len(ms) == 0 {
		return nil, errors.New("no migrations found in the folder")
	}
	for i := 1; i <= len(ms); i++ {
		if _, found := ms[i]; !found {
			return nil, errors.New("migration versions should be an unbroken sequence from one")
		}
	}

	return &ms, nil
}
