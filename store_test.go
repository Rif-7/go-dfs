package main

import (
	"bytes"
	"testing"
)

func TestPathTransformFunc(t *testing.T) {
	key := "testeyname"
	pathKey := CASPathTransformFunc(key)
	expectedOriginalKey := "2b15a221d67f31153ffab04a5a03cba489e98772"
	expectedPathName := "2b15a/221d6/7f311/53ffa/b04a5/a03cb/a489e/98772"

	if pathKey.Original != expectedOriginalKey {
		t.Errorf("have %s want %s", pathKey.Original, expectedOriginalKey)
	}

	if pathKey.PathName != expectedPathName {
		t.Errorf("have %s want %s", pathKey.PathName, expectedPathName)
	}

}

func TestStore(t *testing.T) {
	opts := StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
	}

	s := NewStore(opts)

	data := bytes.NewReader([]byte("test data"))
	if err := s.writeStream("testfile", data); err != nil {
		t.Error(err)
	}

}
