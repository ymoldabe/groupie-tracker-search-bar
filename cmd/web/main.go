package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
)

type Artist struct {
	ID            int            `json:"id"`
	IMAGE         string         `json:"image"`
	NAME          string         `json:"name"`
	MEMBERS       map[int]string `json:"members"`
	CREATION_DATE int            `json:"creationDate"`
	FIRST_ALBUM   string         `json:"firstAlbum"`
	LOCATIONS     string         `json:"locations"`
	CONCERT_DATES string         `json:"concertDates"`
	RELATIONS     string         `json:"relations"`
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", group)
	log.Println("go to run http://localhost:8000/")
	err := http.ListenAndServe(":8000", mux)
	log.Fatal(err)
}

func group(w http.ResponseWriter, r *http.Request) {
	artistDATA := []Artist{}
	par := "https://groupietrackers.herokuapp.com/api/artists"
	mumu, _ := url.Parse(par)
	jsonCHIK, err := url.ParseQuery(mumu.RawQuery)
	if err != nil {
		fmt.Print("sadadsa")
		return
	}
	// defer jsonCHIK.Body.Close()
	fmt.Println(jsonCHIK)
	if r.Method != "GET" {
		Error(w, http.StatusBadRequest)
		return
	}
	if r.URL.Path != "/" {
		Error(w, 404)
		return
	}
	// err := json.Unmarshal([]byte(jsonCHIK), &artistDATA)
	if err != nil {
		Error(w, http.StatusInternalServerError)
		fmt.Println("tutPipec")
		return
	}
	fmt.Println(artistDATA)
}

// id := r.URL.Query().Get("id")
