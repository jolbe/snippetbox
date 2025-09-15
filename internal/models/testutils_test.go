package models

import (
	"database/sql"
	"os"
	"testing"
)

func newTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("mysql", "test_web:asdf@tcp(127.0.0.1:3306)/test_snippetbox?parseTime=true&multiStatements=true")
	if err != nil {
		t.Fatal(err)
	}

	// Execute the setup.sql script to create the database tables and dummy data
	if err := runSQLFile(db, "./testdata/setup.sql"); err != nil {
		t.Fatal(err)
	}

	// Execute the teardown.sql script and close the connection pool
	t.Cleanup(func() {
		if err := runSQLFile(db, "./testdata/teardown.sql"); err != nil {
			t.Fatal(err)
		}
		db.Close()
	})

	return db
}

func runSQLFile(db *sql.DB, filename string) error {
	script, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	_, err = db.Exec(string(script))
	if err != nil {
		return err
	}

	return nil
}
