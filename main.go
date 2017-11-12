package main

import (
	"net/http"
	"log"
	"sync"
	"html/template"
	"path/filepath"
)

// templ represents a single template
type templateHandler struct {
	once sync.Once
	filename string
	templ *template.Template
}

// ServeHttp handles and the HTTP request
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(
			template.ParseFiles(
				filepath.Join("templates", t.filename)))
	})
	t.templ.Execute(w, r)
}

func main()  {
	http.Handle("/", &templateHandler{filename: "chat.html"})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Listen and serve:", err)
	}
}