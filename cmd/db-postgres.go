// PostgreSQL data access.
// AGPL license. Copyright 2024 K Cartlidge.

package main

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

// dbPostgres details a connection to PostgreSQL.
type dbPostgres struct {
	connection *sql.DB
}

// connect connects to the database.
func (d *dbPostgres) connect(connectionString string) error {
	db, err := sql.Open("postgres", connectionString)
	if err == nil {
		err = db.Ping()
		if err == nil {
			d.connection = db
		}
	}
	return err
}

// ensureMigrationsTable creates the `migratable_state` table
// unless it already exists.
func (d *dbPostgres) ensureMigrationsTable() error {
	stmt := `
CREATE TABLE IF NOT EXISTS migratable_state
(
	id             BIGSERIAL                NOT NULL,
	target_version INT,
	direction      CHARACTER VARYING(4)     NOT NULL,
	migration      CHARACTER VARYING(200)   NOT NULL,
	actioned       TIMESTAMP WITH TIME ZONE NOT NULL,
	PRIMARY KEY    (id)
)
`
	return d.exec(stmt)
}

// removeMigrationsTable drops the `migratable_state` table.
func (d *dbPostgres) removeMigrationsTable() error {
	stmt := `DROP TABLE IF EXISTS migratable_state`
	return d.exec(stmt)
}

// getCurrentVersion retrieves the version number (0+).
func (d *dbPostgres) getCurrentVersion() (int, error) {
	stmt := `
SELECT target_version FROM migratable_state
 ORDER BY id DESC LIMIT 1
`
	v := 0
	if err := d.connection.QueryRow(stmt).Scan(&v); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
	}
	return v, nil
}

// execMigration executes the migration and updates the
// `migratable_state` table accordingly in a transaction.
func (d *dbPostgres) execMigration(statement string, setVersion int, name string, direction string) error {
	// No SQL injection risk as parameter is pre-validated.
	name = limitTo(name, 200)
	mtSql := fmt.Sprintf(
		"INSERT INTO migratable_state (direction, target_version, migration, actioned) VALUES ('%s',%d,'%s',NOW())",
		direction, setVersion, name)
	return d.exec(statement, mtSql)
}

// exec executes the given SQL statement(s) in a transaction.
func (d *dbPostgres) exec(statement ...string) error {
	var err error
	var tx *sql.Tx

	// Start a transaction (to allow for multiple statements).
	if tx, err = d.connection.Begin(); err != nil {
		return err
	}

	// Action each statement.
	for i := range statement {
		if _, err = tx.Exec(statement[i]); err != nil {
			pgErr, ok := err.(*pq.Error)
			_ = tx.Rollback()
			if ok {
				return errors.New(pgErr.Error())
			}
			return errors.New("unknown error executing statement")
		}
	}

	// Done.
	err = tx.Commit()
	return err
}
