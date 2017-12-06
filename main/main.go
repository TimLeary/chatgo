package main

import (
	"net/http"
	"log"
	"sync"
	"html/template"
	"path/filepath"
	"os"
	"chatgo/trace"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/gplus"
	"math"
	"github.com/markbates/goth/gothic"
	"github.com/gorilla/sessions"
	"github.com/gorilla/mux"
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

func init() {
	store := sessions.NewFilesystemStore(os.TempDir(), []byte("goth-data"))
	// set the maxLength of the cookies stored on the disk to a larger number to prevent issues with:
	// securecookie: the value is too long
	// when using OpenID Connect , since this can contain a large amount of extra information in the id_token

	// Note, when using the FilesystemStore only the session.ID is written to a browser cookie, so this is explicit for the storage on disk
	store.MaxLength(math.MaxInt64)

	gothic.Store = store
}

func main() {

	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)
	fmt.Println(exPath)

	err = godotenv.Load(exPath + string(filepath.Separator) + ".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var host = os.Getenv("HOST")
	var port = os.Getenv("ADDRESS")

	fmt.Println(os.Getenv("SECURITY_KEY"))
	googleRedirect := "http://"+ host + port + "/auth/callback/gplus"
	goth.UseProviders(
		gplus.New(os.Getenv("GOOGLE_CLIENT_ID"),os.Getenv("GOOGLE_CLIENT_SECRET"), googleRedirect))


	room := newRoom()
	room.tracer = trace.New(os.Stdout)

	router := mux.NewRouter()
	router.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	router.Handle("/login", &templateHandler{filename: "login.html"})

	router.HandleFunc("/auth/login/{provider}", loginHandler)
	router.HandleFunc("/auth/callback/{provider}", callbackHandler)
	router.HandleFunc("auth/logout", logoutHandler)
	router.Handle("/room", room)

	pwd, _ := os.Getwd()
	assets := pwd + "/../assets/"
	fs := http.FileServer(http.Dir( assets ))
	router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets", fs))

	// get the room going
	go room.run()
	// start the web server

	http.Handle("/", router)
	log.Println("Starting web server on", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal("Error ListenAndServe:", err)
	}
}
