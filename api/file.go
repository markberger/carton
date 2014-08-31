package api

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/markberger/carton/common"
	"github.com/markberger/carton/db"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// Most of this function was inspired by Sanat Gersappa's blog post:
// sanatgersappa.blogspot.com/2013/03/handling-multiple-file-uploads-in-go.html

func fileHandler(
	db db.DbManager,
	jar *sessions.CookieStore,
	dest string,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check client has permission to upload a file
		session, _ := jar.Get(r, "carton-session")
		if _, ok := session.Values["user"]; !ok {
			http.Error(w, "No user logged in", http.StatusUnauthorized)
			return
		}

		if r.Method == "GET" {
			files, err := db.GetAllFiles()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			b, err := json.Marshal(files)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			var out bytes.Buffer
			json.Indent(&out, b, "", "\t")
			out.WriteTo(w)
		} else if r.Method == "POST" {
			reader, err := r.MultipartReader()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			c := &common.CartonFile{}
			user, _ := session.Values["user"].(string)
			c.Owner = user
			for {
				part, err := reader.NextPart()
				if err == io.EOF {
					break
				}

				if part.FileName() == "" {
					continue
				}

				if fileExists(dest + part.FileName()) {
					http.Error(w, "File already exists", http.StatusBadRequest)
					return
				}

				filePath := filepath.Join(dest, part.FileName())
				f, err := os.Create(filePath)
				defer f.Close()
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				c.Name = part.FileName()
				c.Path, _ = filepath.Abs(filePath)
				hasher := md5.New()
				writer := io.MultiWriter(f, hasher)
				_, err = io.Copy(writer, part)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				c.Md5Hash = fmt.Sprintf("%x", hasher.Sum(nil))
				c.PwdHash = nil
			}
			db.AddFile(c)
			w.WriteHeader(http.StatusCreated)
			fmt.Fprintln(w, "upload succeeded")
		} else {
			return404(w)
		}
	})
}

func singleFileHandler(db db.DbManager) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		md5Hash := mux.Vars(r)["hash"]
		c, err := db.GetFileByHash(md5Hash)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, "error decoding file")
			return
		}
		if c == nil {
			return404(w)
			return
		}
		if c.PwdHash == nil {
			// When file is downloaded, the file name is c.Name
			w.Header().Set("Content-Disposition", "attachment; filename="+c.Name)
			http.ServeFile(w, r, c.Path)
		} else {
			return404(w)
		}
	})
}

func fileExists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}
	return false
}
