package main

import (
	"database/sql"
	"errors"
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
	version_number INT,
	actioned       TIMESTAMP WITH TIME ZONE NOT NULL,
	PRIMARY KEY    (id)
)
`
	_, err := d.connection.Exec(stmt)
	return err
}

// getCurrentVersion retrieves the version number (0+).
func (d *dbPostgres) getCurrentVersion() (int, error) {
	stmt := `
SELECT version_number FROM migratable_state
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
