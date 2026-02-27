package main

import (
	"net/http"
	"strings"

	"github.com/Bention99/fin-planalyse/internal/auth"
)

func (a *app) handleLoginGet(w http.ResponseWriter, r *http.Request) {
	if err := a.tpl.ExecuteTemplate(w, "login.html", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (a *app) handleLoginPost(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form", http.StatusBadRequest)
		return
	}

	email := strings.TrimSpace(r.FormValue("email"))
	password := r.FormValue("password")

	user, err := a.queries.GetUserByEmail(r.Context(), email)
	if err != nil {
		a.renderLoginWithError(w, "Invalid email or password")
		return
	}

	match, err := auth.CheckPasswordHash(password, user.HashedPassword)
	if err != nil || !match {
		a.renderLoginWithError(w, "Invalid email or password")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    user.ID.String(),
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	})

	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

func (a *app) renderLoginWithError(w http.ResponseWriter, msg string) {
	data := struct {
		Error string
	}{
		Error: msg,
	}

	if err := a.tpl.ExecuteTemplate(w, "login.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}