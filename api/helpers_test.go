package api

import (
	"bytes"
	"github.com/markberger/carton/common"
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
	c := &common.CartonFile{
		"test file 1",
		"md5 hash",
		"/fake/path",
		[]byte("file pass"),
		"owner",
	}
	db.AddFile(c)
	c2, err := db.GetFileByHash("md5 hash")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if c2 != c {
		t.Error("Failed to get file by hash")
	}

	c2 = db.GetFileByName("test file 1")
	if c2 != c {
		t.Error("Failed to get file by name")
	}

	db = NewMockDbManager(true)
	err = db.RegisterUser("test user", []byte("test pass"))
	if err == nil {
		t.Error("Failed to throw error during registration")
	}
}
