package main

import (
	"net/http"
	"os"
	"io"
	"strings"

	"github.com/google/uuid"
)

func (a *app) handleUpload(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}

	file, fh, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "failed to read file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	filename := strings.ToLower(fh.Filename)
	if !strings.HasSuffix(filename, ".pdf") {
		http.Error(w, "only .pdf files are allowed", http.StatusBadRequest)
		return
	}
	if ct := strings.ToLower(fh.Header.Get("Content-Type")); ct != "application/pdf" && ct != "" {
		http.Error(w, "only PDF uploads are allowed", http.StatusBadRequest)
		return
	}

	if err := os.MkdirAll("./uploads", 0o755); err != nil {
		http.Error(w, "cannot create uploads directory", http.StatusInternalServerError)
		return
	}

	dstPath := "./uploads/statement.pdf"
	dst, err := os.Create(dstPath)
	if err != nil {
		http.Error(w, "cannot create file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, "cannot save file", http.StatusInternalServerError)
		return
	}

	userID, ok := userIDFromContext(r.Context())
	if !ok {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}

	uid := uuid.NullUUID{
		UUID:  userID,
		Valid: true,
	}

	catsIncome, err := a.queries.GetCategoriesIncome(r.Context(), uid)
	if err != nil {
		http.Error(w, "failed to load categories: "+err.Error(), http.StatusInternalServerError)
		return
	}

	incomeCategories := make([]string, 0, len(catsIncome))
	for _, row := range catsIncome {
		incomeCategories = append(incomeCategories, row.Name)
	}

	catsExpense, err := a.queries.GetCategoriesExpense(r.Context(), uid)
	if err != nil {
		http.Error(w, "failed to load categories: "+err.Error(), http.StatusInternalServerError)
		return
	}

	expenseCategories := make([]string, 0, len(catsExpense))
	for _, row := range catsExpense {
		expenseCategories = append(expenseCategories, row.Name)
	}

	err = uploadFile(incomeCategories, expenseCategories)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/home#transactions", http.StatusSeeOther)

	//w.Write([]byte("upload successful"))
}