package api

import (
	"code.google.com/p/go.crypto/bcrypt"
	"github.com/gorilla/sessions"
	"net/http"
	"net/http/httptest"
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

func TestLogoutHandler(t *testing.T) {
	jar := sessions.NewCookieStore([]byte("secret key"))
	logoutHandle := logoutHandler(jar)
	test := GenerateHandleTester(t, logoutHandle)

	// Test logout without a logged in user is registered
	// as a bad request.
	w := test("POST", "/logout", url.Values{})
	if w.Code != http.StatusBadRequest {
		t.Errorf(
			"POST /logout: without user returned %v. Expected %v.",
			w.Code,
			http.StatusBadRequest,
		)
	}

	// Test logout with logged in user returns as a good
	// request.
	req, err := http.NewRequest("POST", "/logout", nil)
	if err != nil {
		t.Errorf("%v", err)
	}
	w = httptest.NewRecorder()
	session, _ := jar.Get(req, "carton-session")
	session.Values["user"] = "test user"
	session.Save(req, w)

	logoutHandle.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf(
			"POST /logout: with user returned %v. Expected %v.",
			w.Code,
			http.StatusOK,
		)
	}
}

func TestRegisterHandle(t *testing.T) {
	mockDb := NewMockDbManager(true)
	jar := sessions.NewCookieStore([]byte("secret key"))
	registerHandle := registerHandler(mockDb, jar)
	test := GenerateHandleTester(t, registerHandle)

	// Test GET request
	w := test("GET", "/register", url.Values{})
	if w.Code != http.StatusNotFound {
		t.Errorf(
			"GET /register returned %v. Expected %v",
			w.Code,
			http.StatusNotFound,
		)
	}

	goodParams := url.Values{
		"user":  []string{"test user"},
		"pass1": []string{"test pass"},
		"pass2": []string{"test pass"},
	}

	// Test bad inputs and possible registration error
	badParams := [...]url.Values{
		url.Values{},
		url.Values{
			"user": []string{"test user"},
		},
		url.Values{
			"pass1": []string{"test pass"},
			"pass2": []string{"test pass"},
		},
		// This should fail because we created a mockDb that will
		// throw an error when attempting to register a new user.
		goodParams,
	}

	expectedCode := []int{
		http.StatusBadRequest,
		http.StatusBadRequest,
		http.StatusBadRequest,
		http.StatusInternalServerError,
	}

	for i := range badParams {
		w := test("POST", "/register", badParams[i])
		if w.Code != expectedCode[i] {
			t.Errorf(
				"POST /register: bad input returned %v. Expected %v.",
				w.Code,
				expectedCode[i],
			)
		}
	}

	// Test register fails when user already logged in
	req, err := http.NewRequest("POST", "/register", nil)
	if err != nil {
		t.Errorf("%v", err)
	}
	w = httptest.NewRecorder()
	session, _ := jar.Get(req, "carton-session")
	session.Values["user"] = "test user"
	session.Save(req, w)
	registerHandle.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf(
			"POST /register: when user logged in returned %v. Expected %v.",
			w.Code,
			http.StatusBadRequest,
		)
	}

	// Test that user is successfully registered
	mockDb = NewMockDbManager(false)
	registerHandle = registerHandler(mockDb, jar)
	test = GenerateHandleTester(t, registerHandle)
	w = test("POST", "/register", goodParams)
	if w.Code != http.StatusCreated {
		t.Errorf(
			"POST /register: good input returned %v. Expected %v.",
			w.Code,
			http.StatusCreated,
		)
	}
	if !mockDb.IsUser("test user") {
		t.Error("test user was not added to the database.")
	}
}
