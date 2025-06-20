package main

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const defaultRootFolderName = "storage"

type PathTransformFunc func(string) PathKey

func CASPathTransformFunc(key string) PathKey {
	hash := sha1.Sum([]byte(key))
	hashStr := hex.EncodeToString(hash[:])

	blockSize := 5
	sliceLen := len(hashStr) / 5
	paths := make([]string, sliceLen)

	for i := range sliceLen {
		from, to := i*blockSize, i*blockSize+blockSize
		paths[i] = hashStr[from:to]
	}

	return PathKey{
		PathName: strings.Join(paths, "/"),
		Filename: hashStr,
	}

}

type PathKey struct {
	PathName string
	Filename string
}

func (p PathKey) FullPath() string {
	return fmt.Sprintf("%s/%s", p.PathName, p.Filename)
}

func (p PathKey) FirstPathName() string {
	paths := strings.Split(p.PathName, "/")
	if len(paths) == 0 {
		return ""
	}

	return paths[0]
}

type StoreOpts struct {
	// Name of the root storage folder
	Root              string
	PathTransformFunc PathTransformFunc
}

var DefaultPathTransformFunc = func(key string) PathKey {
	return PathKey{
		PathName: key,
		Filename: key,
	}
}

type Store struct {
	StoreOpts
}

func NewStore(opts StoreOpts) *Store {
	if opts.PathTransformFunc == nil {
		opts.PathTransformFunc = DefaultPathTransformFunc
	}

	if len(opts.Root) == 0 {
		opts.Root = defaultRootFolderName
	}

	return &Store{
		StoreOpts: opts,
	}
}

// Has returns true if the file identified by
// the given ID and key exists on local disk.
func (s *Store) Has(id, key string) bool {
	pathKey := s.PathTransformFunc(key)
	fullPathWithRoot := fmt.Sprintf("%s/%s/%s", s.Root, id, pathKey.FullPath())

	_, err := os.Stat(fullPathWithRoot)
	return !errors.Is(err, os.ErrNotExist)
}

// Clear removes all files
func (s *Store) Clear() error {
	return os.RemoveAll(s.Root)
}

// Delete removes the file from local disk
func (s *Store) Delete(id, key string) error {
	pathKey := s.PathTransformFunc(key)
	fullPath := filepath.Join(s.Root, id, pathKey.FullPath())

	// Delete the file itself
	if err := os.Remove(fullPath); err != nil && !os.IsNotExist(err) {
		return err
	}
	log.Printf("deleted file [%s]", pathKey.Filename)

	// Define the base boundary (e.g., storage)
	boundary := s.Root
	dir := filepath.Dir(fullPath)

	// Recursively remove empty parent directories until the boundary
	for dir != boundary && strings.HasPrefix(dir, boundary) {
		isEmpty, err := isDirEmpty(dir)
		if err != nil {
			return err
		}
		if !isEmpty {
			break
		}
		if err := os.Remove(dir); err != nil {
			return err
		}
		dir = filepath.Dir(dir)
	}

	return nil
}

// isDirEmpty checks if a directory has no entries
func isDirEmpty(path string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer f.Close()

	// Read one entry to check emptiness
	_, err = f.Readdirnames(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err
}

func (s *Store) Read(id, key string) (int64, io.Reader, error) {
	return s.readStream(id, key)
}

func (s *Store) Write(id string, key string, r io.Reader) (int64, error) {
	return s.writeStream(id, key, r)

}

func (s *Store) readStream(id, key string) (int64, io.ReadCloser, error) {
	pathKey := s.PathTransformFunc(key)
	fullPathWithRoot := fmt.Sprintf("%s/%s/%s", s.Root, id, pathKey.FullPath())

	file, err := os.Open(fullPathWithRoot)
	if err != nil {
		return 0, nil, err
	}

	fi, err := file.Stat()
	if err != nil {
		return 0, nil, err
	}

	return fi.Size(), file, nil
}

func (s *Store) OpenFileForWriting(id, key string) (*os.File, error) {
	pathKey := s.PathTransformFunc(key)
	pathNameWithRoot := fmt.Sprintf("%s/%s/%s", s.Root, id, pathKey.PathName)
	if err := os.MkdirAll(pathNameWithRoot, os.ModePerm); err != nil {
		return nil, err
	}

	fullPathWithRoot := fmt.Sprintf("%s/%s/%s", s.Root, id, pathKey.FullPath())

	return os.Create(fullPathWithRoot)
}

func (s *Store) WriteDecrypt(encKey []byte, id string, key string, r io.Reader) (int64, error) {
	f, err := s.OpenFileForWriting(id, key)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	n, err := copyDecrypt(encKey, r, f)
	return int64(n), err
}

func (s *Store) writeStream(id string, key string, r io.Reader) (int64, error) {
	f, err := s.OpenFileForWriting(id, key)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	return io.Copy(f, r)
}
