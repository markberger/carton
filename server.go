package main

import (
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/markberger/carton/api"
	"github.com/markberger/carton/db"
	"net/http"
	"os"
)

// Maps files in path to be served at the given url

func serveDir(m *mux.Router, url string, path string) {
	m.PathPrefix(url).Handler(
		http.StripPrefix(
			url,
			http.FileServer(http.Dir(path)),
		),
	)
}

func serveFile(url string, path string) {
	http.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, path)
	})
}

func registerVendor(m *mux.Router) {
	serveDir(m, "/static/bootstrap/", "./vendor/bootstrap-3.2.0/")
	serveDir(m, "/static/dropzone/", "./vendor/dropzone-3.10.2/")
	serveDir(m, "/static/angular/", "./vendor/angular/")
	serveDir(m, "/", "./public/")
}

func main() {
	b, _ := db.NewBoltManager("./bolt.db")
	m := mux.NewRouter()
	m.StrictSlash(true)
	jar := sessions.NewCookieStore([]byte("secret key"))
	os.Mkdir("./carton_files", os.ModeDir|0764)
	api.RegisterHandlers(m, b, jar, "./carton_files")
	registerVendor(m)
	http.Handle("/", m)

	// ClearHandler is for gorilla/sessions. There will be a
	// memory leak without it
	http.ListenAndServe(":8080", context.ClearHandler(http.DefaultServeMux))
}
