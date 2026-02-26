package main

import (
	"net/http"
	"strings"

	"github.com/Bention99/fin-planalyse/internal/auth"
	"github.com/Bention99/fin-planalyse/internal/database"
)

func (a *app) handleRegisterGet(w http.ResponseWriter, r *http.Request) {
	if err := a.tpl.ExecuteTemplate(w, "register.html", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (a *app) handleRegisterPost(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form", http.StatusBadRequest)
		return
	}

	email := strings.TrimSpace(r.FormValue("email"))
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
}