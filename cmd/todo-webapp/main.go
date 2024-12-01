package main

import (
	"database/sql"
	"fmt"
	"slices"

	"github.com/labstack/echo/v4"
	"github.com/mi-wada/go_playground/todo-webapp/handler"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := initDB()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	e := echo.New()
	e.GET("/healthz", handler.NewHandler(db).GetHealthz)
	e.POST("/tasks", handler.NewHandler(db).PostTasks)
	e.Start(":8080")
}

func initDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "todo-webapp.sqlite")
	if err != nil {
		return nil, err
	}

	createTableQueries := []string{
		`CREATE TABLE IF NOT EXISTS tasks(
			id TEXT NOT NULL,
			content TEXT NOT NULL,
			status TEXT NOT NULL,
			deadline TIMESTAMP,
			created_at TIMESTAMP NOT NULL,
			PRIMARY KEY (id)
		)`,
	}
	for q := range slices.Values(createTableQueries) {
		_, err := db.Exec(q)
		if err != nil {
			return nil, fmt.Errorf("failed to create table: %w", err)
		}
	}
	return db, nil
}
