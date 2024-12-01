package main

import (
	"archive/zip"
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/ikawaha/kagome-dict/ipa"
	"github.com/ikawaha/kagome/v2/tokenizer"
	"golang.org/x/text/encoding/japanese"

	_ "github.com/mattn/go-sqlite3"
)

const (
	AuzoraBaseURL = "https://www.aozora.gr.jp"
)

func main() {
	db, err := initDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	authorURL := "https://www.aozora.gr.jp/index_pages/person879.html"

	entries, err := findEntries(authorURL)
	if err != nil {
		log.Fatal(err)
	}
	for entry := range slices.Values(entries) {
		text, err := downloadText(entry.ZipURL)
		if err != nil {
			log.Fatal(err)
		}

		res, err := db.Exec(`
INSERT OR REPLACE INTO
	authors(author_id, author)
	VALUES (?, ?)
		`, entry.AuthorID, entry.Author)
		if err != nil {
			log.Fatal(err)
		}

		res, err = db.Exec(`
INSERT OR REPLACE INTO
	contents(author_id, title_id, title, content)
	VALUES (?, ?, ?, ?)
		`, entry.AuthorID, entry.TitleID, entry.Title, text)
		if err != nil {
			log.Fatal(err)
		}
		docID, err := res.LastInsertId()
		if err != nil {
			log.Fatal(err)
		}

		t, err := tokenizer.New(ipa.Dict(), tokenizer.OmitBosEos())
		if err != nil {
			log.Fatal(err)
		}

		seg := t.Wakati(text)
		_, err = db.Exec(`
INSERT OR REPLACE INTO
	contents_fts(docid, words)
	VALUES(?, ?)
		`, docID, strings.Join(seg, " "),
		)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Registered: %s", entry.Title)
	}
}

func initDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "database.sqlite")
	if err != nil {
		return nil, err
	}

	createTableQueries := []string{
		`CREATE TABLE IF NOT EXISTS authors(author_id TEXT, author TEXT, PRIMARY KEY
  (author_id))`,
		`CREATE TABLE IF NOT EXISTS contents(author_id TEXT, title_id TEXT, title
  TEXT, content TEXT, PRIMARY KEY (author_id, title_id))`,
		`CREATE VIRTUAL TABLE IF NOT EXISTS contents_fts USING fts4(words)`,
	}
	for q := range slices.Values(createTableQueries) {
		_, err := db.Exec(q)
		if err != nil {
			return nil, err
		}
	}
	return db, nil
}

type Entry struct {
	AuthorID string
	Author   string
	TitleID  string
	Title    string
	BookURL  string
	ZipURL   string
}

func findEntries(authorURL string) ([]Entry, error) {
	doc, err := getDoc(authorURL)
	if err != nil {
		return nil, err
	}

	entries := []Entry{}

	author := doc.Find("body > table > tbody > tr:nth-child(1) > td:nth-child(2) > font").Text()
	authorID, err := extractAuthorID(authorURL)
	if err != nil {
		return nil, err
	}

	doc.Find("ol").First().Find("li > a").Each(func(i int, elem *goquery.Selection) {
		// if i > 1 {
		// 	return
		// }
		title := elem.Text()

		bookURL, exists := elem.Attr("href")
		if !exists {
			log.Println("Not found href", title)
			return
		}

		bookURL, err = toAbsURL(AuzoraBaseURL, bookURL)
		if err != nil {
			log.Println(err, title)
			return
		}

		titleID, err := extractTitleID(bookURL)
		if err != nil {
			log.Println(err, title)
			return
		}

		zipURL, err := getZipURL(bookURL)
		if err != nil {
			log.Println(err, title)
			return
		}
		zipURL, err = toAbsURL(bookURL, zipURL)
		if err != nil {
			log.Println(err, title)
			return
		}

		entries = append(
			entries,
			Entry{
				AuthorID: authorID,
				Author:   author,
				Title:    title,
				TitleID:  titleID,
				BookURL:  bookURL,
				ZipURL:   zipURL,
			},
		)
	})

	return entries, nil
}

func getDoc(URL string) (*goquery.Document, error) {
	res, err := http.Get(URL)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("status code is not 200: %d", res.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	return doc, nil
}

func toAbsURL(baseURL, path string) (string, error) {
	base, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}

	ref, err := url.Parse(path)
	if err != nil {
		return "", err
	}

	abs := base.ResolveReference(ref)
	if err != nil {
		return "", err
	}

	return abs.String(), nil
}

func extractAuthorID(bookListURL string) (string, error) {
	parsedURL, err := url.Parse(bookListURL)
	if err != nil {
		return "", err
	}

	re := regexp.MustCompile(`person(\d+)\.html`)
	matches := re.FindStringSubmatch(parsedURL.Path)
	if len(matches) < 2 {
		return "", fmt.Errorf("no match found")
	}

	id, err := strconv.Atoi(matches[1])
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%06d", id), nil
}

func extractTitleID(bookURL string) (string, error) {
	parsedURL, err := url.Parse(bookURL)
	if err != nil {
		return "", err
	}

	re := regexp.MustCompile(`card(\d+)\.html`)
	matches := re.FindStringSubmatch(parsedURL.Path)
	if len(matches) < 2 {
		return "", fmt.Errorf("no match found")
	}

	return matches[1], nil
}

func getZipURL(bookURL string) (string, error) {
	doc, err := getDoc(bookURL)
	if err != nil {
		return "", err
	}

	zipURL := ""
	doc.Find("body > table.download > tbody > tr > td > a").Each(func(i int, elem *goquery.Selection) {
		URL, _ := elem.Attr("href")
		if strings.HasSuffix(URL, ".zip") && !strings.HasSuffix(URL, "ttz.zip") {
			zipURL = URL
		}
	})

	if zipURL == "" {
		return "", errors.New("No zip URL")
	}

	return zipURL, nil
}

func downloadText(zipURL string) (string, error) {
	resp, err := http.Get(zipURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	r, err := zip.NewReader(bytes.NewReader(b), int64(len(b)))
	if err != nil {
		return "", err
	}

	for _, file := range r.File {
		if path.Ext(file.Name) == ".txt" {
			f, err := file.Open()
			defer f.Close()
			if err != nil {
				return "", err
			}

			b, err := io.ReadAll(f)
			if err != nil {
				return "", err
			}

			b, err = japanese.ShiftJIS.NewDecoder().Bytes(b)
			if err != nil {
				return "", err
			}

			return string(b), nil
		}
	}

	return "", nil
}
