package db

type DbManager interface {
	IsUser(user string) bool
	RegisterUser(user string, hash []byte) error
	GetPwdHash(user string) []byte
}
