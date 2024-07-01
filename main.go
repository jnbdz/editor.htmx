package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("sqlite3", "./editor.db")
	if err != nil {
		log.Fatal(err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)

	createTable()

	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/save", saveHandler)
	http.HandleFunc("/load", loadHandler)

	fmt.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func createTable() {
	query := `
	CREATE TABLE IF NOT EXISTS content (
		id INTEGER PRIMARY KEY,
		text TEXT
	);
	`
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	text := r.FormValue("text")
	_, err := db.Exec("INSERT INTO content (text) VALUES (?)", text)
	if err != nil {
		http.Error(w, "Error saving content", http.StatusInternalServerError)
		return
	}
	_, err = fmt.Fprintln(w, "Content saved")
	if err != nil {
		return
	}
}

func loadHandler(w http.ResponseWriter, r *http.Request) {
	row := db.QueryRow("SELECT text FROM content ORDER BY id DESC LIMIT 1")
	var text string
	err := row.Scan(&text)
	if err != nil {
		if err == sql.ErrNoRows {
			text = ""
		} else {
			http.Error(w, "Error loading content", http.StatusInternalServerError)
			return
		}
	}
	_, err = fmt.Fprintln(w, text)
	if err != nil {
		return
	}
}
