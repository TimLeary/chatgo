package main

import (
	"net/http"
	"log"
	"sync"
	"html/template"
	"path/filepath"
	"flag"
	"os"
	"chatgo/trace"
	"strings"
	"fmt"
	"github.com/joho/godotenv"
)

// templ represents a single template
type templateHandler struct {
	once sync.Once
	filename string
	templ *template.Template
}

// ServeHttp handles and the HTTP request
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r  *http.Request) {
	t.once.Do(func() {
		t.templ =  template.Must(
			template.ParseFiles(
				filepath.Join("../templates", t.filename)))
	})
	t.templ.Execute(w, r)
}

// loginHandler handles the third-party login process.
// format: /auth/{action}/{provider}
func loginHandler(w http.ResponseWriter, r *http.Request)  {
	segs := strings.Split(r.URL.Path, "/")
	action := segs[2]
	provider := segs[3]

	switch action {
	case "login":
		log.Println("TODO handle login for", provider)
	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Auth action %s not supported", action)
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	var addr = flag.String("addr", ":8888", "The addr of the  application.")
	flag.Parse() // parse the flags

	//googleRedirect := "http://localhost" + *addr + "/auth/callback/google"

	r := newRoom()
	r.tracer = trace.New(os.Stdout)
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/room", r)
	pwd, _ := os.Getwd()
	assets := pwd + "/../assets/"
	fs := http.FileServer(http.Dir( assets ))
	http.Handle("/assets/", http.StripPrefix("/assets", fs))
	// get the room going
	go r.run()
	// start the web server
	log.Println("Starting web server on", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
