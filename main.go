package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/rbroggi/chatting/trace"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/github"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/objx"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"text/template"
)

func main() {
	addr := flag.String("addr", ":8080", "The address of the application.")
	verbose := flag.Bool("v", false, "Set verbose mode - tracing active")
	credFilePath := flag.String("cred_file", "credentials.txt", "credentials filename with format provider|app_key|app_secret on each row. providers are google and github")
	flag.Parse()
	r := newRoom()
	if *verbose {
		r.tracer = trace.New(os.Stdout)
	}
	credentials := retrieveCredentials(*credFilePath)
	if len(credentials) < 2 {
		log.Fatalf("Not enough providers in credential file")
	}
	gomniauth.SetSecurityKey("chatting-app-go-blueprints-key")
	gomniauth.WithProviders(
		google.New(credentials["google"].key, credentials["google"].secret, fmt.Sprintf("http://%s/auth/callback/google", *addr)),
		github.New(credentials["github"].key, credentials["github"].secret, fmt.Sprintf("http://%s/auth/callback/github", *addr)),
	)
	//Serve static contents
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	//Endpoint to chat must autheticate first
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/room", r)
	//start the room
	go r.run()
	log.Printf("Starting the we server at address %s.", *addr)
	//start web server
	_ = http.ListenAndServe(*addr, nil)
	log.Fatal("Server mux failed")
}

type cred struct {
	key    string
	secret string
}

func retrieveCredentials(credFilePath string) map[string]cred {

	file, err := os.Open(credFilePath)
	if err != nil {
		log.Fatalf("Could not open credentials file: %s", credFilePath)
	}
	defer file.Close()
	credByProvider := make(map[string]cred)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		record := strings.Split(line, "|")
		if len(record) < 3 {
			continue
		}
		credByProvider[record[0]] = cred{
			key:    record[1],
			secret: record[2],
		}
	}
	return credByProvider
}

//templateHandler handles the parsing of templates
type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

//Implements the http.Handler Interface
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	data := map[string]interface{}{
		"Host": r.Host,
	}
	//Decode value in the "auth" cookie sent by the client
	if authCookie, err := r.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}
	t.templ.Execute(w, data)
}
