package main

import (
	"encoding/json"
	"fmt"
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
	files := "./ui/html/artistData.html"
	tmpl, err := template.ParseFiles(files)
	if err != nil {
		Error(w, http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, artistData)
}

func group(w http.ResponseWriter, r *http.Request) {
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
		"./ui/html/body_home.html",
		"./ui/html/footer_partial.html",
		"./ui/html/front.base.html",
	}
	tmpl, err := template.ParseFiles(files...)
	if err != nil {
		Error(w, http.StatusInternalServerError)
		return
	}
	if r.Method == "GET" {
		err = tmpl.Execute(w, groups)
	} else if r.Method == "POST" {
		find := r.FormValue("search")
		fmt.Println(find)

		date_group := []Data_group{}
		err = json.Unmarshal([]byte(jsonData), &date_group)
		if err != nil {
			Error(w, http.StatusInternalServerError)
			return
		}
		Check_coincidence(date_group, find)

	}
}

func Check_coincidence(data []Data_group, find string) {
	//res := []Coincidence{}
	for _, v := range data {
	}

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
