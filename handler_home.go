package main

import (
	"net/http"

	"github.com/Bention99/fin-planalyse/internal/database"
)

func (a *app) handleHome(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromContext(r.Context())
	if !ok {
		http.Error(w, "missing user in context", http.StatusInternalServerError)
		return
	}

	user, err := a.queries.GetUserByID(r.Context(), userID)
	if err != nil {
		http.Error(w, "failed to load user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	cats, err := a.queries.GetCategories(r.Context())
	if err != nil {
		http.Error(w, "failed to load categories: "+err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		User       database.GetUserByIDRow
		Categories []database.Category
	}{
		User:       user,
		Categories: cats,
	}

	if err := a.tpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}