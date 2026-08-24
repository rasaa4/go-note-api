package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"

	"note-api/handler"
	"note-api/repository"
	"note-api/service"
)

func main() {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/test_db?parseTime=true")
	if err != nil {
		log.Fatal("failed to open database", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("failed to connect to the database", err)
	}
	noteRepo := repository.NewmysqlNoteRepository(db)
	noteService := service.NewNoteService(noteRepo)
	noteHandler := handler.NewNoteHandler(noteService)

	mux := http.NewServeMux()
	mux.HandleFunc("/notes", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			noteHandler.Create(w, r)
		case http.MethodGet:
			noteHandler.GetAll(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/notes/get", func(w http.ResponseWriter, r *http.Request) {
		noteHandler.GetByID(w, r)
	})

	log.Println("server starting on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
