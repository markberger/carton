package common

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
)

type CartonFile struct {
	Name    string
	Md5Hash string
	Path    string
	PwdHash []byte
	Owner   string
}

func (c *CartonFile) GobEncode() ([]byte, error) {
	w := new(bytes.Buffer)
	encoder := gob.NewEncoder(w)
	err := encoder.Encode(c.Name)
	if err != nil {
		return nil, err
	}
	err = encoder.Encode(c.Md5Hash)
	if err != nil {
		return nil, err
	}
	err = encoder.Encode(c.Path)
	if err != nil {
		return nil, err
	}
	err = encoder.Encode(c.PwdHash)
	if err != nil {
		return nil, err
	}
	err = encoder.Encode(c.Owner)
	if err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}

func (c *CartonFile) GobDecode(buf []byte) error {
	r := bytes.NewBuffer(buf)
	decoder := gob.NewDecoder(r)
	err := decoder.Decode(&c.Name)
	if err != nil {
		return err
	}
	err = decoder.Decode(&c.Md5Hash)
	if err != nil {
		return err
	}
	err = decoder.Decode(&c.Path)
	if err != nil {
		return err
	}
	err = decoder.Decode(&c.PwdHash)
	if err != nil {
		return err
	}
	return decoder.Decode(&c.Owner)

}

func (c *CartonFile) MarshalJSON() ([]byte, error) {
	attributes := map[string]string{
		"name":  c.Name,
		"hash":  c.Md5Hash,
		"owner": c.Owner,
	}

	b, err := json.Marshal(attributes)
	if err != nil {
		return nil, err
	}
	return b, nil
}
