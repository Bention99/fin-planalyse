package main

import (
	"net/http"
	"strings"

	"github.com/google/uuid"

	"github.com/Bention99/fin-planalyse/internal/database"
)

func (a *app) handleCreateCategory(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form", http.StatusBadRequest)
		return
	}

	userID, ok := userIDFromContext(r.Context())
	if !ok {
		http.Error(w, "missing user in context", http.StatusInternalServerError)
		return
	}

	uid := uuid.NullUUID{
		UUID:  userID,
		Valid: true,
	}

	name := strings.TrimSpace(r.FormValue("name"))
	typ := strings.TrimSpace(r.FormValue("type"))

	if name == "" || (typ != "income" && typ != "expense") {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	_, err := a.queries.CreateCategory(r.Context(), database.CreateCategoryParams{
		Name: name,
		Type: database.TransactionType(typ),
		UserID: uid,
	})
	if err != nil {
		a.renderHomeWithError(w, r, "Category already exists")
		return
	}

	http.Redirect(w, r, "/home", http.StatusSeeOther)
}