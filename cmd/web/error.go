package main

import (
	"log"
	"net/http"
	"text/template"
)

type ErrorList struct {
	ErrorCode int
	ErrorMess string
}

func Error(w http.ResponseWriter, code int) {
	Err := ErrorList{
		ErrorCode: code,
		ErrorMess: http.StatusText(code),
	}
	w.WriteHeader(code)
	tmpl, err := template.ParseFiles("./ui/html/error.html")
	if err != nil {
		log.Fatal(err)
		return
	}
	err = tmpl.Execute(w, Err)
}

// proverka
