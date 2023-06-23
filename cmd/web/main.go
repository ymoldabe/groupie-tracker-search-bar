package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"text/template"
)

type Artist struct {
	ID            int      `json:"id"`
	IMAGE         string   `json:"image"`
	NAME          string   `json:"name"`
	MEMBERS       []string `json:"members"`
	CREATION_DATE int      `json:"creationDate"`
	FIRST_ALBUM   string   `json:"firstAlbum"`
	LOCATIONS     string   `json:"locations"`
	CONCERT_DATES string   `json:"concertDates"`
	RELATIONS     string   `json:"relations"`
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", group)
	log.Println("go to run http://localhost:8000/")
	err := http.ListenAndServe(":8000", mux)
	log.Fatal(err)
	styles := http.FileServer(http.Dir("/home/student/groupie_treker/ui/static/style.css"))
	mux.Handle("/styles/", http.StripPrefix("/styles/", styles))
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

	artistDATA := []Artist{}

	urlWay := "https://groupietrackers.herokuapp.com/api/artists"

	json4ik, err := http.Get(urlWay)
	if err != nil {
		Error(w, http.StatusInternalServerError)
		return
	}
	defer json4ik.Body.Close()

	body, err := ioutil.ReadAll(json4ik.Body)
	if err != nil {
		Error(w, http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal([]byte(body), &artistDATA)
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
	err = tmpl.Execute(w, artistDATA)
}

// id := r.URL.Query().Get("id")
