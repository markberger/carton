package api

import (
	"code.google.com/p/go.crypto/bcrypt"
	"github.com/gorilla/sessions"
	"github.com/markberger/carton/db"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLoginHandler(t *testing.T) {
	mockDb := db.NewMockDbManager(false)
	jar := sessions.NewCookieStore([]byte("secret key"))
	loginHandle := loginHandler(mockDb, jar)
	test := GenerateHandleTester(t, loginHandle)

	// Test GET request
	w := test("GET", "")
	if w.Code != http.StatusNotFound {
		t.Errorf(
			"GET /login returned %v. Expected %v",
			w.Code,
			http.StatusNotFound,
		)
	}

	goodParams := `{"username":"test user", "password":"test pass"}`

	// Test possible combinations of bad inputs
	badParams := [...]string{
		`{}`,
		`{"username":"test user"}`,
		`{"password":"test pass"}`,
		// Should fail because not in database
		goodParams,
	}

	for _, params := range badParams {
		w := test("POST", params)
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

	w = test("POST", goodParams)
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
	w := test("POST", "")
	if w.Code != http.StatusBadRequest {
		t.Errorf(
			"POST /logout: without user returned %v. Expected %v.",
			w.Code,
			http.StatusBadRequest,
		)
	}

	// Test logout with logged in user returns as a good
	// request.
	req, err := http.NewRequest("POST", "", nil)
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
	mockDb := db.NewMockDbManager(true)
	jar := sessions.NewCookieStore([]byte("secret key"))
	registerHandle := registerHandler(mockDb, jar)
	test := GenerateHandleTester(t, registerHandle)

	// Test GET request
	w := test("GET", "")
	if w.Code != http.StatusNotFound {
		t.Errorf(
			"GET /register returned %v. Expected %v",
			w.Code,
			http.StatusNotFound,
		)
	}

	goodParams := `{
		"username": "test user",
		"password1": "test pass",
		"password2": "test pass"
	}`

	// Test bad inputs and possible registration error
	badParams := [...]string{
		`{}`,
		`{"username": "test user"}`,
		`{"password1": "test pass", "password2": "test pass"}`,
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
		w := test("POST", badParams[i])
		if w.Code != expectedCode[i] {
			t.Errorf(
				"POST /register: bad input returned %v. Expected %v.",
				w.Code,
				expectedCode[i],
			)
		}
	}

	// Test register fails when user already logged in
	req, err := http.NewRequest("POST", "", nil)
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
	mockDb = db.NewMockDbManager(false)
	registerHandle = registerHandler(mockDb, jar)
	test = GenerateHandleTester(t, registerHandle)
	w = test("POST", goodParams)
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

func TestStatusHandle(t *testing.T) {
	jar := sessions.NewCookieStore([]byte("secret key"))
	statusHandle := statusHandler(jar)
	test := GenerateHandleTester(t, statusHandle)

	// Check that status failed without a user logged in
	w := test("GET", "")
	if w.Code != http.StatusForbidden {
		t.Errorf(
			"GET /status returned %v. Expected %v",
			w.Code,
			http.StatusForbidden,
		)
	}

	// Check succeeds when a user is logged in
	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		t.Errorf("%v", err)
	}
	w = httptest.NewRecorder()
	session, _ := jar.Get(req, "carton-session")
	session.Values["user"] = "test user"
	session.Save(req, w)
	statusHandle.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf(
			"GET /status with user returned %v. Expected %v",
			w.Code,
			http.StatusOK,
		)
	}
}
