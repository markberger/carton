package api

import (
	"errors"
	"github.com/markberger/carton/common"
	"net/http"
	"net/http/httptest"
	"net/url"
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
	params url.Values,
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
		params url.Values,
	) *httptest.ResponseRecorder {

		req, err := http.NewRequest(
			method,
			"",
			strings.NewReader(params.Encode()),
		)
		if err != nil {
			t.Errorf("%v", err)
		}
		req.Header.Set(
			"Content-Type",
			"application/x-www-form-urlencoded; param=value",
		)
		w := httptest.NewRecorder()
		handleFunc.ServeHTTP(w, req)
		return w
	}
}

type MockDbManager struct {
	registrationError bool
	users             map[string][]byte
}

func NewMockDbManager(registrationError bool) *MockDbManager {
	db := MockDbManager{}
	db.registrationError = registrationError
	db.users = make(map[string][]byte)
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
	return nil
}

func (db *MockDbManager) GetFileByName(name string) (
	*common.CartonFile,
	error,
) {
	return nil, nil
}
