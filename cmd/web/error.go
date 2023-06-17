package main

import (
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
	tmpl, err := template.ParseFiles("/template/index/error.html")
	if err != nil {
		Error(w, 500)
		return
	}
	err = tmpl.Execute(w, Err)
}
