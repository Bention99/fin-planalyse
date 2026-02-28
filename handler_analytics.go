package main

import (
	"net/http"
	"encoding/json"
	"html/template"
)

func (a *app) handleAnalytics(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromContext(r.Context())
	if !ok {
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}

	rows, err := a.queries.GetBalanceOverTime(r.Context(), userID)
	if err != nil {
		http.Error(w, "failed to load analytics", http.StatusInternalServerError)
		return
	}

	type chartPayload struct {
		Labels []string `json:"labels"`
		Values []int64  `json:"values"`
	}

	payload := chartPayload{}

	for _, row := range rows {
		payload.Labels = append(payload.Labels, row.Date.Format("2006-01-02"))
		payload.Values = append(payload.Values, row.Balance)
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, "failed to encode chart data", http.StatusInternalServerError)
		return
	}

	data := struct {
		ChartData template.JS
	}{
		ChartData: template.JS(payloadJSON),
	}

	if err := a.tpl.ExecuteTemplate(w, "analytics.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}