package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"text/template"
)

func artist(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		Error(w, http.StatusBadRequest)
		return
	}
	checkID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || checkID < 1 {
		Error(w, http.StatusNotFound)
		return
	}
	id := strconv.Itoa(checkID)
	artistData := Artist{}
	jsonData, err := getURL("https://groupietrackers.herokuapp.com/api/artists/" + id)
	if err != nil {
		Error(w, http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal([]byte(jsonData), &artistData)
	if err != nil {
		Error(w, http.StatusInternalServerError)
		return
	}
	jsonData, err = getURL(artistData.RELATIONS)
	if err != nil {
		Error(w, http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal([]byte(jsonData), &artistData)
	if err != nil {
		Error(w, http.StatusInternalServerError)
		return
	}
	files := "/home/student/groupie_treker/ui/html/artistData.html"
	tmpl, err := template.ParseFiles(files)
	if err != nil {
		Error(w, http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, artistData)
}

func group(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		Error(w, http.StatusBadRequest)
		return
	}
	if r.URL.Path != "/" {
		Error(w, 404)
		return
	}

	groups := []Artist{}

	jsonData, err := getURL("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		Error(w, http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal([]byte(jsonData), &groups)
	if err != nil {
		Error(w, http.StatusInternalServerError)
		return
	}
	files := []string{
		"/home/student/groupie_treker/ui/html/body_home.html",
		"/home/student/groupie_treker/ui/html/footer_partial.html",
		"/home/student/groupie_treker/ui/html/front.base.html",
	}
	tmpl, err := template.ParseFiles(files...)
	if err != nil {
		Error(w, http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, groups)
}

func getURL(url string) (js []byte, err error) {
	var errorer []byte
	json4ik, err := http.Get(url)
	if err != nil {
		return errorer, err
	}
	defer json4ik.Body.Close()

	body, err := ioutil.ReadAll(json4ik.Body)
	if err != nil {
		return errorer, err
	}
	return body, err
}
