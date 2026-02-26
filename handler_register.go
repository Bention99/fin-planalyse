package main

import (
	"net/http"

	"github.com/Bention99/fin-planalyse/internal/auth"
	"github.com/Bention99/fin-planalyse/internal/database"
)

func (a *app) handleRegister(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if err := a.tpl.ExecuteTemplate(w, "register.html", nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return

	case http.MethodPost:
		email := r.FormValue("email")
		password := r.FormValue("password")
		if email == "" || password == "" {
			http.Error(w, "missing fields", http.StatusBadRequest)
			return
		}

		hash, err := auth.HashPassword(password)
		if err != nil {
			http.Error(w, "could not hash password", http.StatusInternalServerError)
			return
		}

		_, err = a.queries.CreateUser(r.Context(), database.CreateUserParams{
			Email:          email,
			HashedPassword: hash,
		})
		if err != nil {
			http.Error(w, "could not create user: "+err.Error(), http.StatusBadRequest)
			return
		}

		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}