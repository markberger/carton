package db

import (
	"github.com/markberger/carton/common"
)

type DbManager interface {
	IsUser(user string) bool
	RegisterUser(user string, hash []byte) error
	GetPwdHash(user string) []byte
	AddFile(c *common.CartonFile) error
	GetAllFiles() ([]*common.CartonFile, error)
	GetFileByName(name string) (*common.CartonFile, error)
	Close() error
}
