// Migrations collection.
// AGPL license. Copyright 2024 K Cartlidge.

package main

import (
	"errors"
	"fmt"
	"path"
	"strconv"
	"strings"
)

// migrations holds all loaded migrations.
type migrations map[int]migration

// migration holds one complete migration.
type migration struct {
	Name, Up, Down string
	Display        string
	Version        int
}

// loadMigrations reads in all the migrations in the
// sub-folders of the specified folder.
func loadMigrations(folder string, version int) (*migrations, error) {
	ms := migrations{}

	// Scan for migrations folders.
	if ok, err := exists(folder); !ok || (err != nil) {
		return nil, errors.New("cannot read migrations folder")
	}

	folders, err := getSubFolders(folder)
	if err != nil {
		return nil, err
	}

	for _, f := range folders {
		// Extract the leading version number and the name.
		bits := strings.SplitN(f, " ", 2)
		if len(bits) != 2 {
			return nil, loadingError(f, errors.New("bad migration folder name format"))
		}
		digits := bits[0]
		name := bits[1]
		v, err := strconv.Atoi(digits)
		if err != nil {
			return nil, loadingError(f, errors.New("migration folder bad version format"))
		}
		nm := limitTo(name, 200)

		// Populate the migration.
		subFolder := path.Join(folder, f)
		sqlUp, err := loadFile(subFolder, "up.sql", true)
		if err != nil {
			return nil, loadingError(f, err)
		}
		sqlDn, err := loadFile(subFolder, "down.sql", true)
		if err != nil {
			return nil, loadingError(f, err)
		}
		m := migration{
			Version: v,
			Name:    nm,
			Up:      sqlUp,
			Down:    sqlDn,
		}

		// Gather by version.
		ms[v] = m
	}

	// Check the version numbers are complete and go from 1 upward.
	if len(ms) == 0 {
		return nil, errors.New("empty migrations folder")
	}
	for i := 1; i <= len(ms); i++ {
		if _, found := ms[i]; !found {
			return nil, loadingError(strconv.Itoa(i), errors.New("migration versions should be an unbroken sequence from one"))
		}
	}

	// Now we know the highest version, update the Display values.
	h := ms.getHighest()
	for i := 1; i <= h; i++ {
		ms[i] = ms[i].setDisplay(h)
	}

	return &ms, nil
}

// getHighest returns the highest numbered migration's version (or 0).
func (m *migrations) getHighest() int {
	return len(*m)
}

// getMigration returns the given migration.
// It assumes the index is valid; check before calling.
func (m *migrations) getMigration(i int) migration {
	return (*m)[i]
}

// loadingError wraps an error with a name prefix.
func loadingError(name string, err error) error {
	return errors.New(fmt.Sprintf("`%s`: %s", name, err.Error()))
}

// setDisplay generates the Display field on the migration.
// This cannot be done during initial population as the width
// of the version number is not known until the highest value
// has been gathered.
func (m migration) setDisplay(highest int) migration {
	wid := len(strconv.Itoa(highest))
	m.Display = fmt.Sprintf("%"+strconv.Itoa(wid)+"v", m.Version)
	m.Display = fmt.Sprintf("%s %s", m.Display, m.Name)
	return m
}
