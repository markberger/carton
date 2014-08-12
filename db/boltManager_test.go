package db

import (
	"github.com/boltdb/bolt"
	"io/ioutil"
	"os"
	"testing"
)

// Helper functions

func getTempDb() string {
	tmpDirPath := os.TempDir()
	f, err := ioutil.TempFile(tmpDirPath, "carton_dbTest")
	if err != nil {
		return ""
	}
	f.Close()
	return f.Name()
}

// BoltManager tests

func TestNewBoltManager(t *testing.T) {
	tempDb := getTempDb()
	if tempDb == "" {
		t.Skip("Cannot create temp file")
	}

	m, err := NewBoltManager(tempDb)
	if err != nil {
		t.Errorf("Error when calling NewBoltManager: %v", err)
	}

	createdUsersBucket := false
	m.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("users"))
		if b != nil {
			createdUsersBucket = true
		}
		return nil
	})

	if !createdUsersBucket {
		t.Errorf("'users' bucket does not exist")
	}

	m.Close()
	os.Remove(tempDb)
}

func TestUserRegistration(t *testing.T) {
	tempDb := getTempDb()
	if tempDb == "" {
		t.Skip("Cannot create temp file")
	}

	m, _ := NewBoltManager(tempDb)

	if m.IsUser("fake user") {
		t.Errorf("IsUsers returns true for 'fake user'")
	}
	if m.GetPwdHash("fake user") != nil {
		t.Errorf("GetPwdHash returns non nil value for 'fake user'")
	}

	m.RegisterUser("test user", []byte("test hash"))
	if !m.IsUser("test user") ||
		string(m.GetPwdHash("test user")) != "test hash" {
		t.Errorf("failed to register 'test user'")
	}

	m.Close()
	os.Remove(tempDb)
}
