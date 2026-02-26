package main

import (
	"net/http"

	"github.com/Bention99/fin-planalyse/internal/database"
)

func (a *app) handleCreateCategory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	typ := r.FormValue("type")

	if name == "" || (typ != "income" && typ != "expense") {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	_, err := a.queries.CreateCategory(
		r.Context(),
		database.CreateCategoryParams{
			Name: name,
			Type: database.TransactionType(typ),
		},
	)
	if err != nil {
		http.Error(w, "could not create category: "+err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}