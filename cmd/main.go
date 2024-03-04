// Main entry point.
// AGPL license. Copyright 2024 K Cartlidge.

package main

import (
	"errors"
	"fmt"
	"os"
)

var db idb

func main() {

	// Welcome.
	fmt.Println()
	fmt.Println("MIGRATABLE")
	fmt.Println()

	// Gather the user request.
	c, err := getConfig()
	if err == nil {
		c.describe()

		// Fetch the migrations.
		var m *migrations
		m, err = loadMigrations(c.folder)
		if err == nil {
			fmt.Printf("  Migrations loaded: %v\n", len(*m))

			// Get a database connection.
			connStr := os.Getenv(c.connEnv)
			if len(connStr) == 0 {
				err = errors.New("cannot read the environment variable")
			} else {
				fmt.Println()
				fmt.Println("DATABASE")
				db = &dbPostgres{}
				err = ensureDB(db, connStr)
				if err == nil {

					// Show the current state.
					v, err := db.getCurrentVersion()
					if err == nil {
						fmt.Println("  Database is at migration version", v)
					}
				}
			}
		}
	}

	// Handle any errors.
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
