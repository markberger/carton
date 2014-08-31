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
	GetFileByHash(name string) (*common.CartonFile, error)
	GetFileByName(hash string) *common.CartonFile
	DeleteFile(hash string) error
	Close() error
}
