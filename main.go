package main

import (
	"flag"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/github"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/signature"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"os"
	"github.com/stretchr/objx"
)

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

func (t *templateHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})

	data := map[string]interface{}{
		"Host": req.Host,
	}

	if authCookie, err := req.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}
	t.templ.Execute(res, data)
}

func main() {
	var addr = flag.String("addr", ":8080", "The addr of the application.")
	flag.Parse()

	gomniauth.SetSecurityKey(signature.RandomKey(64))
	gomniauth.WithProviders(
		google.New(os.Getenv("GOOGLE_CLIENT_ID"), os.Getenv("GOOGLE_SECRET"), "http://localhost:8080/auth/callback/google"),
		github.New(os.Getenv("GITHUB_CLIENT_ID"), os.Getenv("GITHUB_SECRET"), "http://localhost:8080/auth/callback/github"),
	)

	chatRoom := newRoom()
	http.Handle("/", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/room", chatRoom)
	go chatRoom.run()

	log.Println("Starting web server on", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
