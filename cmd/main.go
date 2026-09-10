package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"

	"note-api/handler"
	"note-api/repository"
	"note-api/service"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("failed to load .env")
	}
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		dbUser,
		dbPassword,
		dbHost,
		dbPort,
		dbName,
	)
	db, err := sql.Open("mysql", dsn)
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

	userRepo := repository.NewmysqlUserRepo(db)
	userService := service.NewUserService(userRepo)
	authHandler := handler.NewAuthHandler(userService)

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
	mux.HandleFunc("/notes/update", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		noteHandler.Update(w, r)
	})

	mux.HandleFunc("/notes/delete", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		noteHandler.DeleteByID(w, r)
	})

	mux.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method is not allowed", http.StatusMethodNotAllowed)
			return
		}
		authHandler.Register(w, r)

	})
	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method is not allowed", http.StatusMethodNotAllowed)
			return
		}
		authHandler.Login(w, r)

	})

	log.Println("server starting on :8080")
	handler := handler.EnableCORS(mux)

	if err := http.ListenAndServe(":8080", handler); err != nil {

		log.Fatal(err)
	}

}
