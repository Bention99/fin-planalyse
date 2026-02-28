package main

import (
	"net/http"
	"github.com/google/uuid"

	"github.com/Bention99/fin-planalyse/internal/database"
)

type HomeData struct {
	User         database.GetUserByIDRow
	Categories		[]database.GetCategoriesRow
	CategoriesIncome   []database.GetCategoriesIncomeRow
	CategoriesExpense   []database.GetCategoriesExpenseRow
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

	uid := uuid.NullUUID{
		UUID:  userID,
		Valid: true,
	}

	cats, err := a.queries.GetCategories(r.Context(), uid)
	if err != nil {
		http.Error(w, "failed to load categories: "+err.Error(), http.StatusInternalServerError)
		return
	}

	catsIncome, err := a.queries.GetCategoriesIncome(r.Context(), uid)
	if err != nil {
		http.Error(w, "failed to load categories: "+err.Error(), http.StatusInternalServerError)
		return
	}

	catsExpense, err := a.queries.GetCategoriesExpense(r.Context(), uid)
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
		Categories:		cats,
		CategoriesIncome:   catsIncome,
		CategoriesExpense:   catsExpense,
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

	uid := uuid.NullUUID{
		UUID:  userID,
		Valid: true,
	}

	cats, err := a.queries.GetCategories(r.Context(), uid)
	if err != nil {
		http.Error(w, "failed to load categories: "+err.Error(), http.StatusInternalServerError)
		return
	}

	catsIncome, err := a.queries.GetCategoriesIncome(r.Context(), uid)
	if err != nil {
		http.Error(w, "failed to load categories", http.StatusInternalServerError)
		return
	}

	catsExpense, err := a.queries.GetCategoriesExpense(r.Context(), uid)
	if err != nil {
		http.Error(w, "failed to load categories: "+err.Error(), http.StatusInternalServerError)
		return
	}

	txs, err := a.queries.GetTransactions(r.Context(), userID)
	if err != nil {
		http.Error(w, "failed to load transactions", http.StatusInternalServerError)
		return
	}

	data := HomeData{
		User:         user,
		Categories:		cats,
		CategoriesIncome:   catsIncome,
		CategoriesExpense:   catsExpense,
		Transactions: txs,
		Error:        msg,
	}

	if err := a.tpl.ExecuteTemplate(w, "index.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}