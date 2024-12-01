package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := initDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if len(os.Args) < 2 {
		log.Fatal("Please provide a query")
	}

	cmd := os.Args[1]
	switch cmd {
	case "artists":
		rows, err := db.Query(`
			SELECT
				author_id,
				author
			FROM
				authors
		`)
		if err != nil {
			log.Fatal(err)
		}

		for rows.Next() {
			var authorID, author string
			err = rows.Scan(&authorID, &author)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(authorID, author)
		}
	case "titles":
		if len(os.Args) < 3 {
			log.Fatal("Please provide an author ID")
		}
		authorID := os.Args[2]
		rows, err := db.Query(`
			SELECT
				title_id,
				title
			FROM
				contents
			WHERE
				author_id = ?
		`, authorID)
		if err != nil {
			log.Fatal(err)
		}

		for rows.Next() {
			var titleID, title string
			err = rows.Scan(&titleID, &title)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(titleID, title)
		}
	case "content":
		if len(os.Args) < 4 {
			log.Fatal("valid: content <author ID> <title ID>")
		}
		authorID := os.Args[2]
		titleID := os.Args[3]
		rows, err := db.Query(`
			SELECT
				content
			FROM
				contents
			WHERE
				author_id = ? AND
				title_id = ?
		`, authorID, titleID)
		if err != nil {
			log.Fatal(err)
		}

		for rows.Next() {
			var content string
			err = rows.Scan(&content)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(content)
		}
	case "query":
		if len(os.Args) < 3 {
			log.Fatal("Please provide a query")
		}
		query := os.Args[2]
		rows, err := db.Query(`
			SELECT
				a.author,
				c.title
			FROM
				contents c
			INNER JOIN authors a
				ON a.author_id = c.author_id
			INNER JOIN contents_fts f
				ON c.rowid = f.docid
				AND words MATCH ?
		`, query)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		for rows.Next() {
			var author, title string
			err = rows.Scan(&author, &title)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(author, title)
		}
	default:
		help()
	}
}

func initDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "database.sqlite")
	if err != nil {
		return nil, err
	}

	return db, nil
}

func help() {
	fmt.Println("Usage: aozora-searcher [artists|titles|content]")
	os.Exit(1)
}
