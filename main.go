package main

import "log"
import "net/http"
import "path/filepath"
import "sync"
import "text/template"

func main() {
	_ = http.ListenAndServe(":8080", &templateHandler{filename: "chat.html"})
	log.Fatal("Server mux failed")
}

//Implements the http.Handler Interface
type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	t.templ.Execute(w, nil)
}
