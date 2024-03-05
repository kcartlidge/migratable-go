// Interface for database access, plus migration table creation.
// AGPL license. Copyright 2024 K Cartlidge.

package main

import "fmt"

// idb details a connection to a database.
type idb interface {

	// connect connects to the database.
	connect(string) error

	// ensureMigrationsTable creates the `migratable_state` table
	// unless it already exists.
	ensureMigrationsTable() error

	// removeMigrationsTable drops the `migratable_state` table.
	removeMigrationsTable() error

	// getCurrentVersion retrieves the version number (0+).
	getCurrentVersion() (int, error)

	// execMigration executes the migration and updates the
	// `migratable_state` table accordingly in a transaction.
	execMigration(statement string, setVersion int, name string, direction string) error

	// exec executes the given SQL statement(s) in a transaction.
	exec(statement ...string) error
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
