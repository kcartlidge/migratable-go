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
	fmt.Println("MIGRATABLE v1.0.0")
	fmt.Println()

	// Gather the user request.
	c, err := getConfig()
	if err == nil {
		c.describe()

		// Fetch the migrations.
		var m *migrations
		m, err = loadMigrations(c.folder)
		if err == nil {
			fmt.Printf("  Found %v\n", len(*m))

			// Get a database connection.
			connStr := os.Getenv(c.connEnv)
			if len(connStr) == 0 {
				err = errors.New("cannot read the environment variable")
			} else {
				fmt.Println()
				fmt.Println("INFO")
				db = &dbPostgres{}
				err = ensureDB(db, connStr)
				if err == nil {

					// Show the current state.
					var version = 0
					version, err = db.getCurrentVersion()
					if err == nil {
						fmt.Println("  Database is at migration version", version)

						// Work out what to do and do it.
						if c.action != "info" {
							target := version
							switch c.action {
							case "reset":
								target = 0
								break
							case "latest":
								target = m.getHighest()
								break
							case "next":
								if version < m.getHighest() {
									target = version + 1
								}
								break
							case "back":
								if version > 0 {
									target = version - 1
								}
								break
							case "target":
								target = c.target
								break
							}
							fmt.Println("  Targeting migration version", target)

							// Sanity check.
							if target < 0 || target > m.getHighest() {
								err = errors.New(fmt.Sprintf("cannot target version %v", target))
							} else {
								if c.action != "reset" && target == version {
									err = errors.New("already at the requested version")
								} else {
									if version != target {
										fmt.Println("  Migrating from", version, "to", target)
									}
									fmt.Println()
									fmt.Println("APPLYING")

									// Roll forward through all required migrations.
									// Goes from the next to the desired by executing
									// all their UP SQL.
									if target > version {
										for i := version + 1; i <= target; i++ {
											mi := m.getMigration(i)
											fmt.Println(fmt.Sprintf("  UP: %s", mi.Name))
											err = db.execMigration(mi.Up, i, mi.Name, "UP")
											if err != nil {
												break
											}
										}
									}

									// Roll backward through all required migrations.
									// Goes from the current (down) to the desired by
									// executing the DOWN SQL of all but the target.
									if target < version {
										for i := version; i > target; i-- {
											mi := m.getMigration(i)
											fmt.Println(fmt.Sprintf("  DOWN: %s", mi.Name))
											err = db.execMigration(mi.Down, i-1, mi.Name, "DOWN")
											if err != nil {
												break
											}
										}
									}

									// Show the result.
									version, err = db.getCurrentVersion()
									if err == nil {
										fmt.Println("  Database is at migration version", version)
									}

									// Remove the table if a `reset` was requested.
									if c.action == "reset" {
										err = db.removeMigrationsTable()
										if err == nil {
											fmt.Println("  Table `migratable_state` removed")
										}
									}
								}
							}
						}
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
