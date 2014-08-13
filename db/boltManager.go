package db

import (
	"github.com/boltdb/bolt"
	"github.com/markberger/carton/common"
)

type BoltManager struct {
	db *bolt.DB
}

func NewBoltManager(dbPath string) (*BoltManager, error) {
	db, err := bolt.Open(dbPath, 0644, nil)
	if err != nil {
		return nil, err
	}
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("users"))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte("files"))
		return err
	})

	if err != nil {
		return nil, err
	}
	return &BoltManager{db}, nil
}

func (m *BoltManager) Close() error {
	return m.db.Close()
}

func (m *BoltManager) IsUser(user string) bool {
	pwdHash := m.GetPwdHash(user)
	if pwdHash == nil {
		return false
	}
	return true
}

func (m *BoltManager) GetPwdHash(user string) []byte {
	var pwdHash []byte
	m.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("users"))
		v := b.Get([]byte(user))
		if v == nil {
			pwdHash = nil
		} else {
			pwdHash = v
		}
		return nil
	})
	return pwdHash
}

func (m *BoltManager) RegisterUser(user string, hash []byte) error {
	err := m.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("users"))
		err := b.Put([]byte(user), hash)
		return err
	})
	return err
}

func (m *BoltManager) AddFile(c *common.CartonFile) error {
	err := m.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("files"))
		file, err := c.GobEncode()
		if err != nil {
			return err
		}
		err = b.Put([]byte(c.Name), file)
		return err
	})
	return err
}
