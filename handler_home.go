package main

import (
	"net/http"
	"context"

	"github.com/Bention99/fin-planalyse/internal/database"
)

func (a *app) handleHome(w http.ResponseWriter, r *http.Request) {
	cats, err := a.queries.GetCategories(context.Background())
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