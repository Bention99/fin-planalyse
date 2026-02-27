package main

import (
	"net/http"
	"strings"
	"net/mail"
	"fmt"
	"errors"
	"unicode"

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

	_, err := mail.ParseAddress(email)
	if err != nil {
		a.renderRegisterWithError(w, "Please enter a correct Email.")
		return
	}

	err = checkPasswordRequirements(password)
	if err != nil {
		a.renderRegisterWithError(w, fmt.Sprintf("Could not create user: %v", err))
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
		a.renderRegisterWithError(w, "Email already in use. Login instead.")
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (a *app) renderRegisterWithError(w http.ResponseWriter, msg string) {
	data := struct {
		Error string
	}{
		Error: msg,
	}

	if err := a.tpl.ExecuteTemplate(w, "register.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func checkPasswordRequirements(pw string) error {
	if len(pw) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	hasLetter := false
	hasNumber := false
	hasSpecial := false

	for _, r := range pw {
		switch {
		case unicode.IsLetter(r):
			hasLetter = true
		case unicode.IsDigit(r):
			hasNumber = true
		default:
			hasSpecial = true
		}
	}

	if !hasLetter || !hasNumber || !hasSpecial {
		return errors.New("password must contain at least one letter, one number, and one special character")
	}

	return nil
}