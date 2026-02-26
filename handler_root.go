package main

import ("net/http")

func (a *app) handleRoot(w http.ResponseWriter, r *http.Request) {
	_, err := a.getCurrentUserID(r)
	if err != nil {
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/home", http.StatusSeeOther)
}