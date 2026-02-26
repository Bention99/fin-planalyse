package main

import (
	"net/http"

	"github.com/google/uuid"
)

func (a *app) handleDeleteCategory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form", http.StatusBadRequest)
		return
	}

	idStr := r.FormValue("id")
	if idStr == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid uuid", http.StatusBadRequest)
		return
	}

	err = a.queries.DeleteCategory(r.Context(), id)
	if err != nil {
		http.Error(w, "could not delete category: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}