package main

import (
	"github.com/dcso/bloom"
	"crypto/sha1"
	"net/http"
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

type Hit struct {
	Found	bool
}

var filter *bloom.BloomFilter

func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
        w.WriteHeader(status)
}

func check(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		errorHandler(w, r, 404)
	}
	s := []byte(fmt.Sprintf("%X", sha1.Sum([]byte(r.URL.RawQuery))))
	check := filter.Check(s)
	hit := Hit{check}
	js, err := json.Marshal(hit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func checkSHA1(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		errorHandler(w, r, 404)
	}
	s := []byte(r.URL.RawQuery)
	check := filter.Check(s)
	hit := Hit{check}
	js, err := json.Marshal(hit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func main() {
	var err error
	filename := flag.String("f", "pwned-passwords.bloom.gz", "The Bloom filter to load")
	bind := flag.String("b", "127.0.0.1:8000", "The address to which to bind")
	flag.Parse()
	fmt.Printf("Loading Bloom filter from %s...\n", *filename)
	filter, err = bloom.LoadFilter(*filename, true)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(-1)
	}
	fmt.Printf("Listening on %s...", *bind)
	http.HandleFunc("/check", check)
	http.HandleFunc("/check-sha1", checkSHA1)
	fmt.Fprintln(os.Stderr, http.ListenAndServe(*bind, nil))
}

