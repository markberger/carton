package common

import (
	"bytes"
	"testing"
)

func TestFileGob(t *testing.T) {
	c := CartonFile{
		"file name",
		"md5 hash",
		"/fake/path",
		[]byte("file pass"),
		"owner",
	}

	b, err := c.GobEncode()
	if err != nil {
		t.Errorf("Error gob encoding CartonFile: %v", err)
	}
	newFile := CartonFile{}
	err = newFile.GobDecode(b)
	switch {
	case newFile.Name != c.Name:
		t.Error("Names don't match")
	case newFile.Md5Hash != c.Md5Hash:
		t.Error("MD5 hashes don't match")
	case newFile.Path != c.Path:
		t.Error("Paths don't match")
	case !bytes.Equal(newFile.PwdHash, c.PwdHash):
		t.Error("Password hashes don't match")
	case newFile.Owner != c.Owner:
		t.Error("Owners don't match")
	}
}
