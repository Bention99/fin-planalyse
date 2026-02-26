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

	_ = userID

	cats, err := a.queries.GetCategories(r.Context())
	if err != nil {
		http.Error(w, "failed to load categories: "+err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		Categories []database.Category
	}{
		Categories: cats,
	}

	if err := a.tpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}