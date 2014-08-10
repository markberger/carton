package api

import (
	"code.google.com/p/go.crypto/bcrypt"
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/markberger/carton/db"
	"net/http"
)

func return404(w http.ResponseWriter) {
	http.Error(w, "404 page not found", 404)
}

func loginHandler(db db.DbManager, jar *sessions.CookieStore) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			session, _ := jar.Get(r, "carton-session")
			if _, ok := session.Values["user"]; ok {
				http.Error(w, "already signed in", 400)
				return
			}

			username := r.FormValue("user")
			password := r.FormValue("pass")

			if username == "" || password == "" {
				http.Error(w, "bad arguments", 400)
				return
			}

			dbHash := db.GetPwdHash(username)
			if dbHash == nil {
				http.Error(w, "user password combo doesn't exist", 400)
				return
			}

			err := bcrypt.CompareHashAndPassword(dbHash, []byte(password))
			if err != nil {
				http.Error(w, "user password combo doesn't exist", 400)
				return
			}
			session.Values["user"] = username
			session.Save(r, w)
			fmt.Fprintln(w, "login succeeded")
			w.WriteHeader(200)
		} else {
			return404(w)
		}
	})
}

func registerHandler(db db.DbManager, jar *sessions.CookieStore) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			session, _ := jar.Get(r, "carton-session")
			if _, ok := session.Values["user"]; ok {
				http.Error(w, "already signed in", 400)
				return
			}
			username := r.FormValue("user")
			pass1 := r.FormValue("pass1")
			pass2 := r.FormValue("pass2")

			if username == "" ||
				pass1 == "" ||
				pass2 == "" ||
				pass1 != pass2 {
				http.Error(w, "bad arguments", 400)
			}

			if db.IsUser(username) {
				http.Error(w, "user already exists", 400)
				return
			}

			bytePass := []byte(pass1)
			hash, err := bcrypt.GenerateFromPassword(bytePass, bcrypt.DefaultCost)
			if err != nil {
				http.Error(w, "error hashing password", 500)
				return
			}

			err = db.RegisterUser(username, hash)
			if err != nil {
				http.Error(w, "unable to add user", 500)
				return
			}
			session.Values["user"] = username
			session.Save(r, w)
			fmt.Fprintf(w, "Successfully registered %v", username)
			w.WriteHeader(201)
		} else {
			return404(w)
		}
	})
}

func logoutHandler(jar *sessions.CookieStore) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			session, _ := jar.Get(r, "carton-session")
			if _, ok := session.Values["user"]; !ok {
				http.Error(w, "no user to sign out", 400)
				return
			}
			delete(session.Values, "user")
			fmt.Fprintln(w, "Successfully logged out")
			w.WriteHeader(200)
		} else {
			return404(w)
		}
	})
}

func RegisterHandlers(db db.DbManager, jar *sessions.CookieStore) {
	http.Handle("/api/login", loginHandler(db, jar))
	http.Handle("/api/register", registerHandler(db, jar))
	http.Handle("/api/logout", logoutHandler(jar))
}
