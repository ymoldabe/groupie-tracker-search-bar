package main

import (
	"log"
	"net/http"
)

type Artist struct {
	ID                 int                 `json:"id"`
	IMAGE              string              `json:"image"`
	NAME               string              `json:"name"`
	MEMBERS            []string            `json:"members"`
	CREATION_DATE      int                 `json:"creationDate"`
	FIRST_ALBUM        string              `json:"firstAlbum"`
	LOCATIONS          string              `json:"locations"`
	CONCERT_DATES      string              `json:"concertDates"`
	RELATIONS          string              `json:"relations"`
	LOCATION_AND_DATES map[string][]string `json:"datesLocations"`
}

type Artists2 struct {
	ID                 int                 `json:"id"`
	IMAGE              string              `json:"image"`
	NAME               string              `json:"name"`
	MEMBERS            []string            `json:"members"`
	CREATION_DATE      int                 `json:"creationDate"`
	FIRST_ALBUM        string              `json:"firstAlbum"`
	LOCATIONS          string              `json:"locations"`
	CONCERT_DATES      string              `json:"concertDates"`
	RELATIONS          string              `json:"relations"`
	LOCATION_AND_DATES map[string][]string `json:"datesLocations"`
}

type Data_group struct {
	ID                 int      `json:"id"`
	IMAGE              string   `json:"image"`
	NAME               string   `json:"name"`
	MEMBERS            []string `json:"members"`
	LOCATION_AND_DATES Relations
	CREATION_DATE      int    `json:"creationDate"`
	FIRST_ALBUM        string `json:"firstAlbum"`
	RELATIONS          string `json:"relations"`
}

type Relations struct {
	LocationDates map[string][]string `json:"datesLocations"`
}

type Coincidence struct {
	Artist   []Artist
	Artists2 []Artists2
}

func main() {
	mux := http.NewServeMux()
	styles := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static/", styles))
	mux.HandleFunc("/", group)
	mux.HandleFunc("/artist", artist)
	log.Println("Go to run http://localhost:8000/")
	err := http.ListenAndServe(":8000", mux)
	log.Fatal(err)
}
