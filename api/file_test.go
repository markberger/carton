package api

import (
	"github.com/gorilla/sessions"
	"github.com/markberger/carton/db"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestFileUpload(t *testing.T) {
	mockDb := db.NewMockDbManager(false)
	jar := sessions.NewCookieStore([]byte("secret key"))
	tmpDirPath := os.TempDir()
	tmpUploadsPath, err := ioutil.TempDir(tmpDirPath, "cartonUploadTest")
	if err != nil {
		t.Error("Unable to set up tmp directory")
	}
	uploadHandle := fileHandler(mockDb, jar, tmpUploadsPath)
	test := GenerateHandleTester(t, uploadHandle)

	// Check GET request returns 404
	w := test("GET", "")
	if w.Code != http.StatusUnauthorized {
		t.Errorf(
			"GET request returned %v. Expected %v",
			w.Code,
			http.StatusNotFound,
		)
	}

	// Check that someone can't upload a file if they're not logged in
	w = test("POST", "")
	if w.Code != http.StatusUnauthorized {
		t.Errorf(
			"POST request returned %v. Expected %v",
			w.Code,
			http.StatusUnauthorized,
		)
	}

	// Create a tmp file to upload
	tmpFile, err := ioutil.TempFile(tmpDirPath, "cartonUploadTestFile")
	if err != nil {
		t.Errorf("%v", err)
	}
	tmpFile.Close()

	// Test that file upload can succeed
	req, err := newFileUploadRequest("file", tmpFile.Name(), map[string]string{})
	if err != nil {
		t.Errorf("%v", err)
	}
	w = httptest.NewRecorder()
	session, _ := jar.Get(req, "carton-session")
	session.Values["user"] = "test user"
	session.Save(req, w)
	uploadHandle.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Errorf(
			"POST request returned %v. Expected %v",
			w.Code,
			http.StatusCreated,
		)
	}

	c := mockDb.GetFileByName(filepath.Base(tmpFile.Name()))
	if c == nil {
		t.Error("Could not find file")
	}
	switch {
	case c.Name != tmpFile.Name():
	case c.Owner != "test user":
	case c.Path != tmpFile.Name():
		t.Error("Retrieved file does not have expected attributes")
	}

	os.Remove(tmpUploadsPath)
	os.Remove(tmpFile.Name())
}
