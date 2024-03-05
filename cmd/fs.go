// File-system routines.
// AGPL license. Copyright 2024 K Cartlidge.

package main

import (
	"errors"
	"os"
	"path"
	"strings"
)

// exists returns true if the filename can be found.
func exists(filename string) (bool, error) {
	_, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// loadFile returns the file contents.
// Specifying requireContent will return errors if the
// file is empty or all whitespace.
func loadFile(folder string, filename string, requireContent bool) (string, error) {
	sql, err := os.ReadFile(path.Join(folder, filename))
	if err != nil {
		return "", errors.New("unable to read " + filename)
	}
	sql2 := string(sql)
	if len(strings.TrimSpace(sql2)) == 0 {
		return "", errors.New("empty file: " + filename)
	}
	return sql2, nil
}

// getSubFolders returns all sub-folders (no recursion).
func getSubFolders(folder string) ([]string, error) {
	res := []string{}
	files, err := os.ReadDir(folder)
	if err == nil {
		for _, f := range files {
			if f.IsDir() {
				res = append(res, f.Name())
			}
		}
	}
	return res, err
}
