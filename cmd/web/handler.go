package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"
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

	if r.Method == "POST" {
		find := r.FormValue("search")
		fmt.Println(find)
		all_data_group, err := Add_stuckt(w)
		if err != nil {
			Error(w, http.StatusInternalServerError)
			return
		}

		selections := Check_coincidence(w, find, all_data_group)
		err = tmpl.Execute(w, selections)

	} else if r.Method == "GET" {
		err = tmpl.Execute(w, groups)
	}
}

func Check_coincidence(w http.ResponseWriter, find string, all_data_group []Data_group) []Artist {
	res := []Artist{}
	for _, v := range all_data_group {
		for _, j := range v.MEMBERS {
			if strings.Contains(strings.ToLower(j), strings.ToLower(find)) {
				res = append(res, Artist{
					ID:                 v.ID,
					IMAGE:              v.IMAGE,
					NAME:               v.NAME,
					MEMBERS:            v.MEMBERS,
					LOCATION_AND_DATES: v.LOCATION_AND_DATES.LocationDates,
					CREATION_DATE:      v.CREATION_DATE,
					FIRST_ALBUM:        v.FIRST_ALBUM,
					RELATIONS:          v.RELATIONS,
				})
				break
			}
		}
		for key, _ := range v.LOCATION_AND_DATES.LocationDates {
			if strings.Contains(strings.ToLower(key), strings.ToLower(find)) {
				res = append(res, Artist{
					ID:                 v.ID,
					IMAGE:              v.IMAGE,
					NAME:               v.NAME,
					MEMBERS:            v.MEMBERS,
					LOCATION_AND_DATES: v.LOCATION_AND_DATES.LocationDates,
					CREATION_DATE:      v.CREATION_DATE,
					FIRST_ALBUM:        v.FIRST_ALBUM,
					RELATIONS:          v.RELATIONS,
				})
				break
			}
		}
		if strings.Contains(strings.ToLower(v.NAME), strings.ToLower(find)) {
			res = append(res, Artist{
				ID:                 v.ID,
				IMAGE:              v.IMAGE,
				NAME:               v.NAME,
				MEMBERS:            v.MEMBERS,
				LOCATION_AND_DATES: v.LOCATION_AND_DATES.LocationDates,
				CREATION_DATE:      v.CREATION_DATE,
				FIRST_ALBUM:        v.FIRST_ALBUM,
				RELATIONS:          v.RELATIONS,
			})
		} else if strings.Contains((strconv.Itoa(v.CREATION_DATE)), (find)) {
			res = append(res, Artist{
				ID:                 v.ID,
				IMAGE:              v.IMAGE,
				NAME:               v.NAME,
				MEMBERS:            v.MEMBERS,
				LOCATION_AND_DATES: v.LOCATION_AND_DATES.LocationDates,
				CREATION_DATE:      v.CREATION_DATE,
				FIRST_ALBUM:        v.FIRST_ALBUM,
				RELATIONS:          v.RELATIONS,
			})
		} else if strings.Contains((v.FIRST_ALBUM), (find)) {
			res = append(res, Artist{
				ID:                 v.ID,
				IMAGE:              v.IMAGE,
				NAME:               v.NAME,
				MEMBERS:            v.MEMBERS,
				LOCATION_AND_DATES: v.LOCATION_AND_DATES.LocationDates,
				CREATION_DATE:      v.CREATION_DATE,
				FIRST_ALBUM:        v.FIRST_ALBUM,
				RELATIONS:          v.RELATIONS,
			})
		}
	}
	return res
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

func Add_stuckt(w http.ResponseWriter) ([]Data_group, error) {
	res := []Data_group{}
	jsonData, err := getURL("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	var mu sync.Mutex

	var res_stuckt []Data_group
	err = json.Unmarshal(jsonData, &res_stuckt)
	if err != nil {
		return nil, err
	}

	for _, v := range res_stuckt {
		wg.Add(1)
		go func(v Data_group) {
			defer wg.Done()

			jsonData1, err := getURL(v.RELATIONS)
			if err != nil {
				// Обработка ошибок
				return
			}

			var delete Relations
			err = json.Unmarshal([]byte(jsonData1), &delete)
			if err != nil {
				// Обработка ошибок
				return
			}

			mu.Lock()
			res = append(res, Data_group{
				ID:                 v.ID,
				IMAGE:              v.IMAGE,
				NAME:               v.NAME,
				MEMBERS:            v.MEMBERS,
				LOCATION_AND_DATES: delete,
				CREATION_DATE:      v.CREATION_DATE,
				FIRST_ALBUM:        v.FIRST_ALBUM,
				RELATIONS:          v.RELATIONS,
			})
			mu.Unlock()
		}(v)
	}

	wg.Wait()
	return res, nil
}

// func Add_stuckt(w http.ResponseWriter) ([]Data_group, error) {
// 	res := []Data_group{}
// 	delete := Relations{}
// 	res_stuckt := []Data_group{}
// 	jsonData, err := getURL("https://groupietrackers.herokuapp.com/api/artists")
// 	if err != nil {
// 		return []Data_group{}, err
// 	}
// 	err = json.Unmarshal(jsonData, &res_stuckt)
// 	if err != nil {
// 		return []Data_group{}, err
// 	}
// 	for _, v := range res_stuckt {
// 		jsonData1, err := getURL(v.RELATIONS)
// 		if err != nil {
// 			return []Data_group{}, err
// 		}
// 		err = json.Unmarshal([]byte(jsonData1), &delete)
// 		if err != nil {
// 			return []Data_group{}, err
// 		}
// 		res = append(res, Data_group{
// 			NAME:               v.NAME,
// 			MEMBERS:            v.MEMBERS,
// 			LOCATION_AND_DATES: delete,
// 			CREATION_DATE:      v.CREATION_DATE,
// 			FIRST_ALBUM:        v.FIRST_ALBUM,
// 			RELATIONS:          v.RELATIONS,
// 		})
// 		delete = Relations{}

// 	}
// 	return res, nil
// }
