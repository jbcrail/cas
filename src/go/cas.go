package main

import (
	"crypto/sha1"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type Filesystem struct {
	cwd string
}

func (fs *Filesystem) Setcwd(path string) {
	fs.cwd = path + "/"
}

func (fs *Filesystem) Exists(filename string) bool {
	_, err := os.Stat(fs.cwd + filename)
	if err == nil {
		return true
	}
	return os.IsNotExist(err)
}

func (fs *Filesystem) ReadAll(filename string) ([]byte, error) {
	return ioutil.ReadFile(fs.cwd + filename)
}

func (fs *Filesystem) WriteAll(filename string, data []byte) error {
	return ioutil.WriteFile(fs.cwd+filename, data, 0644)
}

func Hash(data []byte) string {
	h := sha1.New()
	h.Write(data)
	return fmt.Sprintf("%x", h.Sum(nil))
}

var fs Filesystem

func RetrieveHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["id"]
	if !fs.Exists(id) {
		http.Error(w, "SHA-1 "+id+" does not exist", http.StatusNotFound)
		return
	}
	content, _ := fs.ReadAll(id)
	contentId := Hash(content)
	if id != contentId {
		http.Error(w, "Content is corrupted", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(content)
}

func StoreHandler(w http.ResponseWriter, req *http.Request) {
	if req.ContentLength < 1 || req.ContentLength > 64*1024*1024 {
		http.Error(w, "Content is less than 1 byte or greater than 64MiB", http.StatusBadRequest)
		return
	}
	content, _ := ioutil.ReadAll(req.Body)
	id := Hash(content)
	if fs.Exists(id) {
		http.Error(w, "", http.StatusConflict)
		return
	}
	w.WriteHeader(http.StatusCreated)
	if err := fs.WriteAll(id, content); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	pwd, _ := os.Getwd()
	fs.Setcwd(pwd)

	port := flag.String("port", "4567", "port")

	r := mux.NewRouter()
	r.HandleFunc("/{id}", RetrieveHandler).Methods("GET")
	r.HandleFunc("/", StoreHandler).Methods("POST")
	http.Handle("/", r)
	http.ListenAndServe(":"+*port, nil)
}
