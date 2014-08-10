package main

import (
	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"github.com/markberger/carton/api"
	"github.com/markberger/carton/db"
	"net/http"
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
	api.RegisterHandlers(b, jar)
	registerVendor()

	// ClearHandler is for gorilla/sessions. There will be a
	// memory leak without it
	http.ListenAndServe(":8080", context.ClearHandler(http.DefaultServeMux))
}
