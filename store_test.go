package main

import (
	"bytes"
	"fmt"
	"io"
	"testing"
)

func TestPathTransformFunc(t *testing.T) {
	key := "testeyname"
	pathKey := CASPathTransformFunc(key)
	expectedFilename := "2b15a221d67f31153ffab04a5a03cba489e98772"
	expectedPathName := "2b15a/221d6/7f311/53ffa/b04a5/a03cb/a489e/98772"

	if pathKey.Filename != expectedFilename {
		t.Errorf("have %s want %s", pathKey.Filename, expectedFilename)
	}

	if pathKey.PathName != expectedPathName {
		t.Errorf("have %s want %s", pathKey.PathName, expectedPathName)
	}

}

func TestStore(t *testing.T) {
	s := newStore()
	defer teardown(t, s)

	for i := range 50 {

		key := fmt.Sprintf("baz_%d", i)
		data := []byte("test data")
		if err := s.writeStream(key, bytes.NewReader(data)); err != nil {
			t.Error(err)
		}

		if ok := s.Has(key); !ok {
			t.Errorf("expected to have key %s", key)
		}

		r, err := s.Read(key)
		if err != nil {
			t.Error(err)
		}

		b, _ := io.ReadAll(r)

		fmt.Println(string(b))

		if string(b) != string(data) {
			t.Errorf("want %s have %s", data, b)
		}

		fmt.Println(string(b))

		if err := s.Delete(key); err != nil {
			t.Error(err)
		}

		if ok := s.Has(key); ok {
			t.Errorf("expected to NOT have key %s", key)
		}

	}

}

func newStore() *Store {
	opts := StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
	}

	return NewStore(opts)
}

func teardown(t *testing.T, s *Store) {
	if err := s.Clear(); err != nil {
		t.Error(err)
	}
}
