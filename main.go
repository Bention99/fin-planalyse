package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"

	"github.com/Bention99/fin-planalyse/internal/database"
)

type app struct {
	db      *sql.DB
	queries *database.Queries
	tpl     *template.Template
}

func main() {
	godotenv.Load(".env")
	dbURL := os.Getenv("DATABASE_URL")
	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	a := &app{
		db:      db,
		queries: database.New(db),
		tpl:     template.Must(template.ParseFiles("templates/index.html")),
	}

	mux := http.NewServeMux()

	mux.Handle("/templates/", http.StripPrefix("/templates/", http.FileServer(http.Dir("./templates"))))

	mux.HandleFunc("/", a.handleHome)
	mux.HandleFunc("/categories", a.handleCreateCategory)
	mux.HandleFunc("/categories/delete", a.handleDeleteCategory)
	addr := ":8080"
	log.Printf("http://localhost%s\n", addr)

	log.Fatal(http.ListenAndServe(":8080", mux))
}