package main

import (
	"database/sql"
	"errors"
	"fmt"
	"html/template"
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
	http.HandleFunc("/list-notes", listNotesHandler)
	http.HandleFunc("/load-note", loadNoteHandler) // Register the loadNoteHandler

	fmt.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func createTable() {
	query := `
	CREATE TABLE IF NOT EXISTS content (
		id INTEGER PRIMARY KEY,
		title TEXT,
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
	title := "Note Title" // Default title or you can handle it differently
	_, err := db.Exec("INSERT INTO content (title, text) VALUES (?, ?)", title, text)
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
		if errors.Is(err, sql.ErrNoRows) {
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

func listNotesHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, title FROM content ORDER BY id DESC")
	if err != nil {
		http.Error(w, "Error listing notes", http.StatusInternalServerError)
		return
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			return
		}
	}(rows)

	var notes []struct {
		ID    int
		Title string
	}
	for rows.Next() {
		var id int
		var title string
		if err := rows.Scan(&id, &title); err != nil {
			http.Error(w, "Error scanning notes", http.StatusInternalServerError)
			return
		}
		notes = append(notes, struct {
			ID    int
			Title string
		}{id, title})
	}
	tmpl := `
	<ul>
		{{range .}}
			<li><a href="#" hx-get="/load-note?id={{.ID}}" hx-target="#editor-container">{{.Title}}</a></li>
		{{end}}
	</ul>
	`
	t := template.Must(template.New("list").Parse(tmpl))
	err = t.Execute(w, notes)
	if err != nil {
		return
	}
}

func loadNoteHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	row := db.QueryRow("SELECT text FROM content WHERE id = ?", id)
	var text string
	err := row.Scan(&text)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Note not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error loading note", http.StatusInternalServerError)
		}
		return
	}
	_, err = fmt.Fprintln(w, text)
	if err != nil {
		return
	}
}
