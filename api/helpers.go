package api

import (
	"errors"
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

type HandleTester func(method string, params url.Values) *httptest.ResponseRecorder

// Given a path and a http.Handler, generate a HandleTester which
// will test its given input against the supplied path and handler.

func GenerateHandleTester(
	t *testing.T,
	path string,
	handleFunc http.Handler,
) HandleTester {

	// Given a method type ("GET", "POST", etc) and parameters,
	// serve the response against the handler and return the
	// ResponseRecorder.

	return func(method string, params url.Values) *httptest.ResponseRecorder {
		req, err := http.NewRequest(
			method,
			path,
			strings.NewReader(params.Encode()),
		)
		if err != nil {
			t.Errorf("%v\n", err)
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
