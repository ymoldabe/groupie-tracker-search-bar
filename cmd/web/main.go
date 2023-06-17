package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", group)
	log.Println("go to run http://localhost:8000/")
	err := http.ListenAndServe(":8000", mux)
	log.Fatal(err)
}

func group(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		Error(w, 404)
		return
	}
	w.Write([]byte("ooo"))
}
