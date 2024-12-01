package testutils

import (
	"database/sql"
	"fmt"
	"slices"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}

	createTableQueries := []string{
		`CREATE TABLE IF NOT EXISTS tasks(id TEXT, content TEXT, status TEXT, deadline TIMESTAMP, created_at TIMESTAMP, PRIMARY KEY (id))`,
	}
	for q := range slices.Values(createTableQueries) {
		_, err := db.Exec(q)
		if err != nil {
			return nil, fmt.Errorf("failed to create table: %w", err)
		}
	}
	return db, nil
}
