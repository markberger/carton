package main

import (
	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"github.com/markberger/carton/api"
	"github.com/markberger/carton/db"
	"net/http"
	"os"
)

// Maps files in path to be served at the given url

func registerDir(url string, path string) {
	http.Handle(
		url,
		http.StripPrefix(
			url,
			http.FileServer(http.Dir(path)),
		),
	)
}

func registerVendor() {
	registerDir("/static/bootstrap/", "./vendor/bootstrap-3.2.0/")
	registerDir("/static/dropzone/", "./vendor/dropzone-3.10.2/")
}

func main() {
	b, _ := db.NewBoltManager("./bolt.db")
	jar := sessions.NewCookieStore([]byte("secret key"))
	os.Mkdir("./carton_files", os.ModeDir|0666)
	api.RegisterHandlers(b, jar, "./carton_files")
	registerVendor()

	// ClearHandler is for gorilla/sessions. There will be a
	// memory leak without it
	http.ListenAndServe(":8080", context.ClearHandler(http.DefaultServeMux))
}
