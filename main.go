package main

import (
	"flag"
	"github.com/rbroggi/chatting/trace"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"text/template"
)

func main() {
	addr := flag.String("addr", ":8080", "The address of the application.")
	verbose := flag.Bool("v", false, "Set verbose mode - tracing active")
	flag.Parse()
	r := newRoom()
	if *verbose {
		r.tracer = trace.New(os.Stdout)
	}
	http.Handle("/", &templateHandler{filename: "chat.html"})
	//Serve static contents
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.Handle("/room", r)
	//start the room
	go r.run()
	log.Printf("Starting the we server at address %s.", *addr)
	//start web server
	_ = http.ListenAndServe(*addr, nil)
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
	t.templ.Execute(w, r)
}
