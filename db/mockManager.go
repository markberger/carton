package db

import (
	"errors"
	"github.com/markberger/carton/common"
)

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
	if _, ok := db.files[c.Md5Hash]; ok {
		return errors.New("File already exists")
	} else {
		db.files[c.Md5Hash] = c
		return nil
	}
}

func (db *MockDbManager) GetFileByName(name string) *common.CartonFile {
	for _, c := range db.files {
		if c.Name == name {
			return c
		}
	}
	return nil
}

func (db *MockDbManager) GetFileByHash(hash string) (
	*common.CartonFile,
	error,
) {
	if _, ok := db.files[hash]; !ok {
		return nil, errors.New("File does not exist")
	} else {
		return db.files[hash], nil
	}
}

func (db *MockDbManager) GetAllFiles() ([]*common.CartonFile, error) {
	files := []*common.CartonFile{}
	for _, v := range db.files {
		files = append(files, v)
	}
	return files, nil
}
