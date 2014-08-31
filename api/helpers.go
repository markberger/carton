package api

import (
	"bytes"
	"errors"
	"github.com/markberger/carton/common"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

/*
 *  API Helpers
 */

func return404(w http.ResponseWriter) {
	http.Error(w, "404 page not found", 404)
}

/*
 *  API Testing Helpers
 */

type HandleTester func(
	method string,
	params string,
) *httptest.ResponseRecorder

// Given the current test runner and an http.Handler, generate a
// HandleTester which will test its the handler against its input

func GenerateHandleTester(
	t *testing.T,
	handleFunc http.Handler,
) HandleTester {

	// Given a method type ("GET", "POST", etc) and
	// parameters, serve the response against the handler
	// and return the ResponseRecorder.

	return func(
		method string,
		params string,
	) *httptest.ResponseRecorder {

		req, err := http.NewRequest(
			method,
			"",
			strings.NewReader(params),
		)
		if err != nil {
			t.Errorf("%v", err)
		}
		req.Header.Set(
			"Content-Type",
			"application/json",
		)
		req.Body.Close()
		w := httptest.NewRecorder()
		handleFunc.ServeHTTP(w, req)
		return w
	}
}

// From Matt Aimonetti's blog post:
// matt.aimonetti.net/posts/2013/07/01/golang-multipart-file-upload-example/
// Creates a new file upload http request with optional extra params
func newFileUploadRequest(
	paramName string,
	path string,
	params map[string]string,
) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, filepath.Base(path))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", "", body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", writer.FormDataContentType())
	return req, nil
}

type MockDbManager struct {
	registrationError bool
	users             map[string][]byte
	files             map[string]*common.CartonFile
}

func NewMockDbManager(registrationError bool) *MockDbManager {
	db := MockDbManager{}
	db.registrationError = registrationError
	db.users = make(map[string][]byte)
	db.files = make(map[string]*common.CartonFile)
	return &db
}

func (db *MockDbManager) Close() error {
	return nil
}

func (db *MockDbManager) IsUser(user string) bool {
	_, ok := db.users[user]
	return ok
}

func (db *MockDbManager) RegisterUser(user string, hash []byte) error {
	if db.registrationError {
		return errors.New("Registration error")
	}
	db.users[user] = hash
	return nil
}

func (db *MockDbManager) GetPwdHash(user string) []byte {
	val, ok := db.users[user]
	if !ok {
		return nil
	} else {
		return val
	}
}

func (db *MockDbManager) AddFile(c *common.CartonFile) error {
	if _, ok := db.files[c.Name]; ok {
		return errors.New("File already exists")
	} else {
		db.files[c.Name] = c
		return nil
	}
}

func (db *MockDbManager) GetFileByName(name string) (
	*common.CartonFile,
	error,
) {
	if _, ok := db.files[name]; !ok {
		return nil, errors.New("File does not exist")
	} else {
		return db.files[name], nil
	}
}

func (db *MockDbManager) GetFileByHash(hash string) *common.CartonFile {
	return nil
}

func (db *MockDbManager) GetAllFiles() ([]*common.CartonFile, error) {
	files := []*common.CartonFile{}
	for _, v := range db.files {
		files = append(files, v)
	}
	return files, nil
}
