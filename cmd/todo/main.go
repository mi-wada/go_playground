package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: todo <command> [arguments]")
		return
	}

	db, err := sql.Open("sqlite3", "./todo.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	createTable(db)

	switch os.Args[1] {
	case "ls":
		listTasks(db)
	case "add":
		if len(os.Args) < 3 {
			fmt.Println("Usage: todo add \"task\"")
			return
		}
		addTask(db, strings.Join(os.Args[2:], " "))
	case "done":
		if len(os.Args) < 3 {
			fmt.Println("Usage: todo done \"task\"")
			return
		}
		markTaskDone(db, strings.Join(os.Args[2:], " "))
	case "rm":
		if len(os.Args) < 3 {
			fmt.Println("Usage: todo rm \"task\"")
			return
		}
		removeTask(db, strings.Join(os.Args[2:], " "))
	case "clear":
		clearTasks(db)
	default:
		fmt.Println("Unknown command:", os.Args[1])
	}
}

func createTable(db *sql.DB) {
	query := `
	CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		task TEXT,
		done BOOLEAN
	);
	`
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}

func listTasks(db *sql.DB) {
	rows, err := db.Query("SELECT id, task, done FROM tasks")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var task string
		var done bool
		err = rows.Scan(&id, &task, &done)
		if err != nil {
			log.Fatal(err)
		}
		status := " "
		if done {
			status = "x"
		}
		fmt.Printf("[%s] %d: %s\n", status, id, task)
	}
}

func addTask(db *sql.DB, task string) {
	_, err := db.Exec("INSERT INTO tasks (task, done) VALUES (?, ?)", task, false)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Added task:", task)
}

func markTaskDone(db *sql.DB, task string) {
	_, err := db.Exec("UPDATE tasks SET done = ? WHERE task = ?", true, task)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Marked task as done:", task)
}

func removeTask(db *sql.DB, task string) {
	_, err := db.Exec("DELETE FROM tasks WHERE task = ?", task)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Removed task:", task)
}

func clearTasks(db *sql.DB) {
	_, err := db.Exec("DELETE FROM tasks")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Cleared all tasks")
}
