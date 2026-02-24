package main

import (
	"html/template"
	"log"
	"net/http"
)

var tpl = template.Must(template.ParseFiles("static/index.html"))

func main() {
	mux := http.NewServeMux()

	// Static assets (css/js/images) at /static/
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	// Homepage rendered by template
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := struct{ Name string }{Name: "Benoit"}
		if err := tpl.Execute(w, data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	addr := ":8080"
	log.Printf("http://localhost%s\n", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}