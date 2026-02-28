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
	//uploadFile()
	godotenv.Load(".env")
	dbURL := os.Getenv("DATABASE_URL")
	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	funcMap := template.FuncMap{
		"formatCents": formatCents,
	}

	a := &app{
		db:      db,
		queries: database.New(db),
		tpl:     template.Must(template.New("index.html").Funcs(funcMap).ParseGlob("templates/*.html")),
	}

	mux := http.NewServeMux()

	fs := http.StripPrefix("/static/", http.FileServer(http.Dir("./static")))
	mux.Handle("GET /static/", fs)
	mux.Handle("HEAD /static/", fs)

	mux.HandleFunc("GET /", a.handleRoot)

	mux.HandleFunc("GET /home", a.requireAuth(a.handleHome))

	mux.HandleFunc("POST /categories", a.requireAuth(a.handleCreateCategory))
	mux.HandleFunc("POST /categories/delete", a.requireAuth(a.handleDeleteCategory))

	mux.HandleFunc("GET /register", a.handleRegisterGet)
	mux.HandleFunc("POST /register", a.handleRegisterPost)

	mux.HandleFunc("GET /login", a.handleLoginGet)
	mux.HandleFunc("POST /login", a.handleLoginPost)

	mux.HandleFunc("POST /logout", a.handleLogout)

	mux.HandleFunc("POST /transactions", a.requireAuth(a.handleCreateTransaction))
	mux.HandleFunc("POST /transactions/delete", a.requireAuth(a.handleDeleteTransaction))

	mux.Handle("POST /upload", a.requireAuth(http.HandlerFunc(a.handleUpload)))

	mux.HandleFunc("GET /analytics", a.requireAuth(a.handleAnalytics))

	addr := ":8080"
	log.Printf("http://localhost%s\n", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}