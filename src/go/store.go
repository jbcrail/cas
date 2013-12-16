package main

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
)

// Errors introduced by storage.
var (
	ErrIdNotExist     = errors.New("storage: id does not exist")
	ErrContentCorrupt = errors.New("storage: corrupted content")
	ErrContentExist   = errors.New("storage: content already exists")
)

type Storage struct {
	root string // Directory for storage content
	mu   sync.RWMutex
}

func (s *Storage) SetRoot(directory string) {
	s.root = directory + "/"
}

func (s *Storage) hash(data []byte) string {
	h := sha1.New()
	h.Write(data)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func (s *Storage) exists(filename string) bool {
	_, err := os.Stat(s.root + filename)
	if err == nil {
		return true
	}
	return os.IsExist(err)
}

// Get content from storage atomically
func (s *Storage) Get(id string) ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if !s.exists(id) {
		return nil, ErrIdNotExist
	}

	content, err := ioutil.ReadFile(s.root + id)
	if err != nil {
		return nil, err
	}

	contentId := s.hash(content)
	if id != contentId {
		return nil, ErrContentCorrupt
	}

	return content, nil
}

// Set content to storage atomically
func (s *Storage) Set(data []byte) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := s.hash(data)

	if s.exists(id) {
		return "", ErrContentExist
	}

	if err := ioutil.WriteFile(s.root+id, data, 0644); err != nil {
		return "", err
	}

	return id, nil
}
