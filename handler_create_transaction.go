package main

import (
	"net/http"
	"strings"
	"time"
	"errors"
	"strconv"

	"github.com/google/uuid"

	"github.com/Bention99/fin-planalyse/internal/database"
)

func (a *app) handleCreateTransaction(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromContext(r.Context())
	if !ok {
		http.Error(w, "missing user in context", http.StatusInternalServerError)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form: "+err.Error(), http.StatusBadRequest)
		return
	}

	categoryInput := strings.TrimSpace(r.FormValue("category_name"))
	parts := strings.Split(categoryInput, " (")
	if len(parts) != 2 {
		return
	}
	name := parts[0]

	cID, err := a.queries.GetCategoryID(r.Context(), name)
	if err != nil {
		http.Error(w, "couldn't fetch category ID", http.StatusBadRequest)
		return
	}

	typePart := strings.TrimSuffix(parts[1], ")")

	aStr := strings.TrimSpace(r.FormValue("amount"))

	amount, err := strconv.ParseFloat(aStr, 64)
	if err != nil {
		http.Error(w, "Invalid amount", http.StatusBadRequest)
		return
	}

	correctedAmount := amount

	if typePart == "expense" {
		if amount >= 0 {
			correctedAmount = amount * -1
		}
	}

	if typePart == "income" {
		if amount <= 0 {
			correctedAmount = amount * -1
		}
	}

	categoryIDStr := cID.String()
	dateStr := strings.TrimSpace(r.FormValue("date"))
	amountStr := strconv.FormatFloat(correctedAmount, 'f', 2, 64)
	isOptional := r.FormValue("is_optional") != ""

	categoryID, err := uuid.Parse(categoryIDStr)
	if err != nil {
		http.Error(w, "invalid category_id", http.StatusBadRequest)
		return
	}

	dt, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		http.Error(w, "invalid date (expected YYYY-MM-DD)", http.StatusBadRequest)
		return
	}
	sqlAmount, err := parseAmountForSQLC(amountStr)
	if err != nil {
		http.Error(w, "invalid amount", http.StatusBadRequest)
		return
	}

	_, err = a.queries.CreateTransaction(r.Context(), database.CreateTransactionParams{
		UserID:     userID,
		CategoryID: categoryID,
		Date:       dt,
		Amount:     sqlAmount,
		IsOptional: isOptional,
	})
	if err != nil {
		http.Error(w, "could not create transaction: "+err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/home#transactions", http.StatusSeeOther)
}

func parseAmountForSQLC(s string) (int64, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, errors.New("amount required")
	}

	s = strings.ReplaceAll(s, ",", ".")

	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, errors.New("invalid amount format")
	}

	cents := int64(f * 100)

	return cents, nil
}