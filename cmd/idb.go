package main

import "fmt"

// idb details a connection to a database.
type idb interface {

	// connect connects to the database.
	connect(string) error

	// ensureMigrationsTable creates the `migratable_state` table
	// unless it already exists.
	ensureMigrationsTable() error

	// getCurrentVersion retrieves the version number (0+).
	getCurrentVersion() (int, error)
}

// ensureDB checks it can connect to the database and
// creates the migrations table if it is not there.
func ensureDB(db idb, connStr string) error {
	err := db.connect(connStr)
	if err == nil {
		fmt.Println("  Connected to the database")
		err = db.ensureMigrationsTable()
		if err == nil {
			fmt.Println("  Table migratable_state exists")
		}
	}
	return err
}
