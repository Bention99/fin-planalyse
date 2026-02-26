package main

import (
	"net/http"

	"github.com/Bention99/fin-planalyse/internal/auth"
)

func (a *app) handleLogin(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if err := a.tpl.ExecuteTemplate(w, "login.html", nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return

	case http.MethodPost:
		email := r.FormValue("email")
		password := r.FormValue("password")

		user, err := a.queries.GetUserByEmail(r.Context(), email)
		if err != nil {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}

		match, err := auth.CheckPasswordHash(password, user.HashedPassword)
		if err != nil || !match {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "session",
			Value:    user.ID.String(),
			HttpOnly: true,
			Path:     "/",
			SameSite: http.SameSiteLaxMode,
		})

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}