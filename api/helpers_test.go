package api

import (
	"bytes"
	"testing"
)

func TestMockDb(t *testing.T) {
	db := NewMockDbManager(false)
	err := db.RegisterUser("test user", []byte("test pass"))
	if err != nil {
		t.Error("Returned unexpected error during registration")
	}
	if !db.IsUser("test user") {
		t.Error("Failed to add user to the mock db")
	}
	hash := db.GetPwdHash("test user")
	if !bytes.Equal(hash, []byte("test pass")) {
		t.Error("Failed to store password in the mock db")
	}

	db = NewMockDbManager(true)
	err = db.RegisterUser("test user", []byte("test pass"))
	if err == nil {
		t.Error("Failed to throw error during registration")
	}
}
