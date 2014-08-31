package db

import (
	"github.com/boltdb/bolt"
	"github.com/markberger/carton/common"
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

func TestAddFile(t *testing.T) {
	tempDb := getTempDb()
	if tempDb == "" {
		t.Skip("Cannot create temp file")
	}

	m, _ := NewBoltManager(tempDb)
	c := &common.CartonFile{
		"file name",
		"md5 hash",
		"/fake/path",
		[]byte("file pass"),
		"owner",
	}

	err := m.AddFile(c)
	if err != nil {
		t.Errorf("Error adding file: %v", err)
	}

	f, err := m.GetFileByName("file name")
	if err != nil {
		t.Errorf("Error getting file: %v", err)
	}
	if f.Name != c.Name {
		t.Error("File names don't match")
	}

	f = m.GetFileByHash("md5 hash")
	if f.Name != c.Name {
		t.Error("File names don't match")
	}

	m.Close()
	os.Remove(tempDb)
}

func TestGetAllFiles(t *testing.T) {
	tempDb := getTempDb()
	if tempDb == "" {
		t.Skip("Cannot create temp file")
	}

	m, _ := NewBoltManager(tempDb)
	c1 := &common.CartonFile{
		"test file 1",
		"md5 hash",
		"/fake/path",
		[]byte("file pass"),
		"owner",
	}
	c2 := &common.CartonFile{
		"test file 2",
		"md5 hash",
		"/fake/path",
		[]byte("file pass"),
		"owner",
	}
	m.AddFile(c1)
	m.AddFile(c2)

	files, err := m.GetAllFiles()
	if err != nil {
		t.Errorf("%v", err)
	}
	if len(files) != 2 {
		t.Errorf("Expected 2 files. Recieved %v", len(files))
	}

	m.Close()
	os.Remove(tempDb)
}
