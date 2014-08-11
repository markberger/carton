package api

import (
	"code.google.com/p/go.crypto/bcrypt"
	"github.com/gorilla/sessions"
	"net/http"
	"net/url"
	"testing"
)

func TestLoginHandler(t *testing.T) {
	mockDb := NewMockDbManager(false)
	jar := sessions.NewCookieStore([]byte("secret key"))
	loginHandle := loginHandler(mockDb, jar)
	test := GenerateHandleTester(t, loginHandle)

	// Test GET request
	w := test("GET", "/login", url.Values{})
	if w.Code != http.StatusNotFound {
		t.Errorf(
			"GET /login returned %v. Expected %v",
			w.Code,
			http.StatusNotFound,
		)
	}

	// Test possible combinations of bad inputs
	badParams := [...]url.Values{
		url.Values{},
		url.Values{
			"user": []string{"test user"},
		},
		url.Values{
			"pass": []string{"test pass"},
		},
		url.Values{
			"user": []string{"test user"},
			"pass": []string{"test pass"},
		},
	}

	for _, params := range badParams {
		w := test("POST", "/login", params)
		if w.Code != http.StatusBadRequest {
			t.Errorf(
				"POST /login: bad input returned %v. Expected %v.",
				w.Code,
				http.StatusBadRequest,
			)
		}
	}

	// Test with good params
	hash, _ := bcrypt.GenerateFromPassword(
		[]byte("test pass"),
		bcrypt.DefaultCost,
	)
	mockDb.RegisterUser("test user", hash)
	goodParams := url.Values{
		"user": []string{"test user"},
		"pass": []string{"test pass"},
	}

	w = test("POST", "/login", goodParams)
	if w.Code != http.StatusOK {
		t.Errorf(
			"POST /login: good input returned %v. Expected %v.",
			w.Code,
			http.StatusOK,
		)
	}
}
