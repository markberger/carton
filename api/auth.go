package api

import (
	"code.google.com/p/go.crypto/bcrypt"
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/markberger/carton/db"
	"net/http"
)

func loginHandler(db db.DbManager, jar *sessions.CookieStore) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			session, _ := jar.Get(r, "carton-session")
			if _, ok := session.Values["user"]; ok {
				http.Error(w, "already signed in", http.StatusBadRequest)
				return
			}

			username := r.PostFormValue("user")
			password := r.PostFormValue("pass")

			if username == "" || password == "" {
				http.Error(w, "bad arguments", http.StatusBadRequest)
				return
			}

			dbHash := db.GetPwdHash(username)
			if dbHash == nil {
				http.Error(
					w,
					"user password combo doesn't exist",
					http.StatusBadRequest,
				)
				return
			}

			err := bcrypt.CompareHashAndPassword(dbHash, []byte(password))
			if err != nil {
				http.Error(
					w,
					"user password combo doesn't exist",
					http.StatusBadRequest,
				)
				return
			}
			session.Values["user"] = username
			session.Save(r, w)
			fmt.Fprintln(w, "login succeeded")
			w.WriteHeader(http.StatusOK)
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
				http.Error(w, "already signed in", http.StatusBadRequest)
				return
			}
			username := r.FormValue("user")
			pass1 := r.FormValue("pass1")
			pass2 := r.FormValue("pass2")

			if username == "" ||
				pass1 == "" ||
				pass2 == "" ||
				pass1 != pass2 {
				http.Error(w, "bad arguments", http.StatusBadRequest)
			}

			if db.IsUser(username) {
				http.Error(w, "user already exists", http.StatusBadRequest)
				return
			}

			bytePass := []byte(pass1)
			hash, err := bcrypt.GenerateFromPassword(bytePass, bcrypt.DefaultCost)
			if err != nil {
				http.Error(
					w,
					"error hashing password",
					http.StatusInternalServerError,
				)
				return
			}

			err = db.RegisterUser(username, hash)
			if err != nil {
				http.Error(
					w,
					"unable to add user",
					http.StatusInternalServerError,
				)
				return
			}
			session.Values["user"] = username
			session.Save(r, w)
			w.WriteHeader(http.StatusCreated)
			fmt.Fprintf(w, "Successfully registered %v", username)
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
				http.Error(w, "no user to sign out", http.StatusBadRequest)
				return
			}
			delete(session.Values, "user")
			session.Save(r, w)
			fmt.Fprintln(w, "Successfully logged out")
			w.WriteHeader(http.StatusOK)
		} else {
			return404(w)
		}
	})
}

func statusHandler(jar *sessions.CookieStore) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			session, _ := jar.Get(r, "carton-session")
			if _, ok := session.Values["user"]; ok {
				w.WriteHeader(200)
				fmt.Fprintln(w, "User is logged in")
				http.Error(w, "No user is signed in", http.StatusForbidden)
			} else {
				http.Error(w, "No user is signed in", http.StatusForbidden)
			}
		} else {
			return404(w)
		}
	})
}

func RegisterHandlers(db db.DbManager, jar *sessions.CookieStore, dest string) {
	http.Handle("/api/login", loginHandler(db, jar))
	http.Handle("/api/register", registerHandler(db, jar))
	http.Handle("/api/logout", logoutHandler(jar))
	http.Handle("/api/status", statusHandler(jar))
	http.Handle("/api/files", fileHandler(db, jar, dest))
}
