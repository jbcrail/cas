package main

import (
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"os"
)

type Storage struct {
	cwd string
}

func (s *Storage) Setcwd(path string) {
	s.cwd = path + "/"
}

func (s *Storage) Exists(filename string) bool {
	_, err := os.Stat(s.cwd + filename)
	if err == nil {
		return true
	}
	return os.IsExist(err)
}

func (s *Storage) ReadAll(filename string) ([]byte, error) {
	return ioutil.ReadFile(s.cwd + filename)
}

func (s *Storage) WriteAll(filename string, data []byte) error {
	return ioutil.WriteFile(s.cwd+filename, data, 0644)
}

func Hash(data []byte) string {
	h := sha1.New()
	h.Write(data)
	return fmt.Sprintf("%x", h.Sum(nil))
}
