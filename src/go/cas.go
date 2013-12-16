package main

import (
	"flag"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

var store Storage

func RetrieveHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["id"]
	content, err := store.Get(id)
	switch err {
	case ErrIdNotExist:
		http.Error(w, "SHA-1 "+id+" does not exist", http.StatusNotFound)
		return
	case ErrContentCorrupt:
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
	id, err := store.Set(content)
	if err == ErrContentExist {
		http.Error(w, "", http.StatusConflict)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(id))
}

func main() {
	pwd, _ := os.Getwd()

	port := flag.String("port", "4567", "port")
	home := flag.String("dir", pwd, "storage directory")
	flag.Parse()

	// set directory for content storage
	store.SetRoot(*home)

	r := mux.NewRouter()
	r.HandleFunc("/{id}", RetrieveHandler).Methods("GET")
	r.HandleFunc("/", StoreHandler).Methods("POST")
	http.Handle("/", r)
	http.ListenAndServe(":"+*port, nil)
}
