package main

import (
	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"github.com/markberger/carton/api"
	"github.com/markberger/carton/db"
	"net/http"
)

func main() {
	b, _ := db.NewBoltManager("./bolt.db")
	jar := sessions.NewCookieStore([]byte("secret key"))
	api.RegisterHandlers(b, jar)

	// ClearHandler is for gorilla/sessions. There will be a
	// memory leak without it
	http.ListenAndServe(":8080", context.ClearHandler(http.DefaultServeMux))
}
