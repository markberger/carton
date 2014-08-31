package api

import (
	"code.google.com/p/go.crypto/bcrypt"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/markberger/carton/db"
	"net/http"
)

type User struct {
	Username string
	Password string
}

func loginHandler(db db.DbManager, jar *sessions.CookieStore) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			session, _ := jar.Get(r, "carton-session")
			if _, ok := session.Values["user"]; ok {
				http.Error(w, "already signed in", http.StatusBadRequest)
				return
			}

			decoder := json.NewDecoder(r.Body)
			var user User
			err := decoder.Decode(&user)
			if err != nil {
				http.Error(w, "error decoding json", http.StatusBadRequest)
				return
			}

			if user.Username == "" || user.Password == "" {
				http.Error(w, "bad arguments", http.StatusBadRequest)
				return
			}

			dbHash := db.GetPwdHash(user.Username)
			if dbHash == nil {
				http.Error(
					w,
					"user password combo doesn't exist",
					http.StatusBadRequest,
				)
				return
			}

			err = bcrypt.CompareHashAndPassword(dbHash, []byte(user.Password))
			if err != nil {
				http.Error(
					w,
					"user password combo doesn't exist",
					http.StatusBadRequest,
				)
				return
			}
			session.Values["user"] = user.Username
			session.Save(r, w)
			// Sets return code to 200
			fmt.Fprintln(w, "login succeeded")
		} else {
			return404(w)
		}
	})
}

type NewUser struct {
	Username  string
	Password1 string
	Password2 string
}

func registerHandler(db db.DbManager, jar *sessions.CookieStore) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			session, _ := jar.Get(r, "carton-session")
			if _, ok := session.Values["user"]; ok {
				http.Error(w, "already signed in", http.StatusBadRequest)
				return
			}

			decoder := json.NewDecoder(r.Body)
			var user NewUser
			err := decoder.Decode(&user)
			if err != nil {
				http.Error(w, "error decoding json", http.StatusBadRequest)
				return
			}

			if user.Username == "" ||
				user.Password1 == "" ||
				user.Password2 == "" ||
				user.Password1 != user.Password2 {
				http.Error(w, "bad arguments", http.StatusBadRequest)
				return
			}

			if db.IsUser(user.Username) {
				http.Error(w, "user already exists", http.StatusBadRequest)
				return
			}

			bytePass := []byte(user.Password1)
			hash, err := bcrypt.GenerateFromPassword(bytePass, bcrypt.DefaultCost)
			if err != nil {
				http.Error(
					w,
					"error hashing password",
					http.StatusInternalServerError,
				)
				return
			}

			err = db.RegisterUser(user.Username, hash)
			if err != nil {
				http.Error(
					w,
					"unable to add user",
					http.StatusInternalServerError,
				)
				return
			}
			session.Values["user"] = user.Username
			session.Save(r, w)
			w.WriteHeader(http.StatusCreated)
			fmt.Fprintf(w, "Successfully registered %v", user.Username)
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
			// Sets return code to 200
			fmt.Fprintln(w, "Successfully logged out")
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
				// Sets return code to 200
				fmt.Fprintln(w, "User is logged in")
			} else {
				http.Error(w, "No user is signed in", http.StatusForbidden)
			}
		} else {
			return404(w)
		}
	})
}

func RegisterHandlers(
	m *mux.Router,
	db db.DbManager,
	jar *sessions.CookieStore,
	dest string,
) {
	m.Handle("/api/auth/login", loginHandler(db, jar))
	m.Handle("/api/auth/register", registerHandler(db, jar))
	m.Handle("/api/auth/logout", logoutHandler(jar))
	m.Handle("/api/auth/status", statusHandler(jar))
	m.Handle("/api/files", fileHandler(db, jar, dest))
	m.Handle("/api/files/{hash}", singleFileHandler(db))
}
