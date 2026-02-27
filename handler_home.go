package main

import (
	"net/http"

	"github.com/Bention99/fin-planalyse/internal/database"
)

type HomeData struct {
	User         database.GetUserByIDRow
	Categories   []database.Category
	Transactions []database.GetTransactionsRow
	Error        string
}

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

	txs, err := a.queries.GetTransactions(r.Context(), userID)
	if err != nil {
		http.Error(w, "failed to load transactions: "+err.Error(), http.StatusInternalServerError)
		return
	}

	data := HomeData{
		User:         user,
		Categories:   cats,
		Transactions: txs,
		Error:        "",
	}

	if err := a.tpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (a *app) renderHomeWithError(w http.ResponseWriter, r *http.Request, msg string) {
	userID, ok := userIDFromContext(r.Context())
	if !ok {
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}

	user, err := a.queries.GetUserByID(r.Context(), userID)
	if err != nil {
		http.Error(w, "failed to load user", http.StatusInternalServerError)
		return
	}

	cats, err := a.queries.GetCategories(r.Context())
	if err != nil {
		http.Error(w, "failed to load categories", http.StatusInternalServerError)
		return
	}

	txs, err := a.queries.GetTransactions(r.Context(), userID)
	if err != nil {
		http.Error(w, "failed to load transactions", http.StatusInternalServerError)
		return
	}

	data := HomeData{
		User:         user,
		Categories:   cats,
		Transactions: txs,
		Error:        msg,
	}

	if err := a.tpl.ExecuteTemplate(w, "index.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}