package main

import (
	"net/http"

	"github.com/Bention99/fin-planalyse/internal/database"
	"github.com/google/uuid"
)

func (a *app) handleDeleteTransaction(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromContext(r.Context())
	if !ok {
		http.Redirect(w, r, "/register", http.StatusSeeOther)
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

	txID, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid uuid", http.StatusBadRequest)
		return
	}

	n, err := a.queries.DeleteTransaction(r.Context(), database.DeleteTransactionParams{
		ID:     txID,
		UserID: userID,
	})
	if err != nil {
		http.Error(w, "could not delete transaction: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if n == 0 {
		http.Error(w, "transaction not found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, "/home", http.StatusSeeOther)
}