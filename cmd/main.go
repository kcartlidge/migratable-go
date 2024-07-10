// Main entry point.
// AGPL license. Copyright 2024 K Cartlidge.

package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var db idb

func main() {

	// Welcome.
	fmt.Println()
	fmt.Println("MIGRATABLE v2.0.2")
	fmt.Println("Postgres Migrations Tool")
	fmt.Println()

	// Gather the user request.
	c, err := getConfig()
	if err == nil {
		fmt.Println()
		c.describe()

		// Get a database connection.
		connStr := os.Getenv(c.connEnv)
		if len(connStr) == 0 {
			err = errors.New("cannot read the environment variable")
		} else {
			db = &dbPostgres{}
			err = ensureDB(db, connStr)
			if err == nil {

				// Get the current version.
				var version = 0
				version, err = db.getCurrentVersion()
				if err == nil {

					// Fetch the migrations.
					var m *migrations
					m, err = loadMigrations(c.folder, version)
					if err == nil {
						fmt.Printf("Migrations loaded (x%v)\n", len(*m))

						// Show the current state.
						displayStatus(m, version)

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

							// Sanity check.
							if target < 0 || target > m.getHighest() {
								err = errors.New(fmt.Sprintf("cannot target version %v", target))
							} else {
								if c.action != "reset" && target == version {
									err = errors.New("already at the requested version")
								} else {
									if version != target {
										section(fmt.Sprintf("MIGRATING FROM %v TO %v", version, target))
									}

									// Roll forward through all required migrations.
									// Goes from the next to the desired by executing
									// all their UP SQL.
									if target > version {
										for i := version + 1; i <= target; i++ {
											mi := m.getMigration(i)
											fmt.Println(fmt.Sprintf("  Migrating UP: %s", mi.Display))
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
											fmt.Println(fmt.Sprintf("  Migrating DOWN: %s", mi.Display))
											err = db.execMigration(mi.Down, i-1, mi.Name, "DOWN")
											if err != nil {
												break
											}
										}
									}

									if err == nil {
										// Show the result.
										version, err = db.getCurrentVersion()
										if err == nil {
											// Remove the table if a `reset` was requested.
											if c.action == "reset" {
												err = db.removeMigrationsTable()
												if err == nil {
													fmt.Println()
													fmt.Println("Table `migratable_state` removed")
												}
											} else {
												displayStatus(m, version)
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
	}

	// Handle any errors.
	if err != nil {
		fmt.Println()
		fmt.Println("ERROR")
		fmt.Println(err.Error())
		fmt.Println()
	} else {
		fmt.Println()
		fmt.Println("Done")
		fmt.Println()
	}
}

// displayStatus shows the migrations with a marker
// for where the current version is.
func displayStatus(m *migrations, version int) {
	section("MIGRATION STATUS")
	h := m.getHighest()
	w := len(strconv.Itoa(h))
	here := strings.Repeat("-", w+1) + "> *"
	shown := false
	if version < 1 {
		fmt.Println(here)
	}
	for i := 1; i <= h; i++ {
		om := (*m)[i]
		fmt.Println(fmt.Sprintf("  %s", om.Display))
		if om.Version == version {
			shown = true
			fmt.Println(here)
		}
	}
	if !shown {
		fmt.Println()
		fmt.Println("Are these the correct migrations for this database?")
	}
}
