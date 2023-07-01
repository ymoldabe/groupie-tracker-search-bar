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

var jsonData1 []byte

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

	// groups := []Artist{}
	// groups_ADG := Coincidence{
	// 	Artist:     groups,
	// 	Data_group: all_data_group,
	// }

	jsonData1, err := getURL("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		Error(w, http.StatusInternalServerError)
		return
	}
	all_data_group, err := Add_stuckt(w, jsonData1)
	// fmt.Println(all_data_group[0])

	// err = json.Unmarshal([]byte(jsonData1), &groups)
	// if err != nil {
	// 	Error(w, http.StatusInternalServerError)
	// 	return
	// }

	if r.Method == "GET" {

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
		err = tmpl.Execute(w, all_data_group)
	} else if r.Method == "POST" {
		find := r.FormValue("search")
		fmt.Println(find)
		if err != nil {
			Error(w, http.StatusInternalServerError)
			return
		}

		selections := Check_coincidence(w, find, all_data_group)
		// selections2 := Coincidence{
		// 	Artist:   groups,
		// 	Artist: selections,
		// }
		tmpl, err := template.ParseFiles("./ui/html/search.html")
		if err != nil {
			Error(w, http.StatusInternalServerError)
			return
		}
		err = tmpl.Execute(w, selections)

	}
}

func Check_coincidence(w http.ResponseWriter, find string, all_data_group []Artist) []Artist {
	res := []Artist{}
	flag := false
	for _, v := range all_data_group {
		if strings.Contains(strings.ToLower(v.NAME), strings.ToLower(find)) ||
			strings.Contains(strconv.Itoa(v.CREATION_DATE), find) ||
			strings.Contains(v.FIRST_ALBUM, find) {
			// fmt.Println(v.LOCATION_AND_DATES.LocationDates)
			res = append(res, Artist{
				ID:                      v.ID,
				IMAGE:                   v.IMAGE,
				NAME:                    v.NAME,
				MEMBERS:                 v.MEMBERS,
				LOCATION_AND_DATES_link: v.LOCATION_AND_DATES.LocationDates,
				CREATION_DATE:           v.CREATION_DATE,
				FIRST_ALBUM:             v.FIRST_ALBUM,
				RELATIONS:               v.RELATIONS,
			})
			continue
		}

		for _, j := range v.MEMBERS {
			if strings.Contains(strings.ToLower(j), strings.ToLower(find)) {
				flag = true
				res = append(res, Artist{
					ID:                      v.ID,
					IMAGE:                   v.IMAGE,
					NAME:                    v.NAME,
					MEMBERS:                 v.MEMBERS,
					LOCATION_AND_DATES_link: v.LOCATION_AND_DATES.LocationDates,
					CREATION_DATE:           v.CREATION_DATE,
					FIRST_ALBUM:             v.FIRST_ALBUM,
					RELATIONS:               v.RELATIONS,
				})
				break
			}
		}
		if flag {
			flag = false
			continue
		}
		for key := range v.LOCATION_AND_DATES.LocationDates {
			if strings.Contains(strings.ToLower(key), strings.ToLower(find)) {
				flag = true
				// fmt.Println(v)
				res = append(res, Artist{
					ID:                      v.ID,
					IMAGE:                   v.IMAGE,
					NAME:                    v.NAME,
					MEMBERS:                 v.MEMBERS,
					LOCATION_AND_DATES_link: v.LOCATION_AND_DATES.LocationDates,
					CREATION_DATE:           v.CREATION_DATE,
					FIRST_ALBUM:             v.FIRST_ALBUM,
					RELATIONS:               v.RELATIONS,
				})
				break
			}
		}
		if flag {
			flag = false
			continue
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

func Add_stuckt(w http.ResponseWriter, jsonData1 []byte) ([]Artist, error) {
	res := []Artist{}

	var wg sync.WaitGroup
	var mu sync.Mutex

	var res_stuckt []Artist
	err := json.Unmarshal(jsonData1, &res_stuckt)
	// fmt.Println(jsonData1)
	if err != nil {
		return nil, err
	}

	for i, v := range res_stuckt {
		wg.Add(1)
		go func(v Artist) {
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
			res = append(res, Artist{
				ID:                      v.ID,
				IMAGE:                   v.IMAGE,
				NAME:                    v.NAME,
				MEMBERS:                 v.MEMBERS,
				LOCATION_AND_DATES:      delete,
				CREATION_DATE:           v.CREATION_DATE,
				FIRST_ALBUM:             v.FIRST_ALBUM,
				RELATIONS:               v.RELATIONS,
				LOCATION_AND_DATES_link: v.LOCATION_AND_DATES.LocationDates,
			})
			fmt.Println(res[0].LOCATION_AND_DATES)
			res = append(res[i].LOCATION_AND_DATES_link, res[i].LOCATION_AND_DATES)

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

// func Add_stuckt(w http.ResponseWriter) ([]Data_group, error) {
// 	res := []Data_group{}

// 	jsonData, err := getURL("https://groupietrackers.herokuapp.com/api/artists")
// 	if err != nil {
// 		return nil, err
// 	}

// 	var res_stuckt []Data_group
// 	err = json.Unmarshal(jsonData, &res_stuckt)
// 	if err != nil {
// 		return nil, err
// 	}

// 	for _, v := range res_stuckt {
// 		jsonData1, err := getURL(v.RELATIONS)
// 		if err != nil {
// 			// Handle errors
// 			return nil, err
// 		}

// 		var delete Relations
// 		err = json.Unmarshal([]byte(jsonData1), &delete)
// 		if err != nil {
// 			// Handle errors
// 			return nil, err
// 		}

// 		res = append(res, Data_group{
// 			ID:                 v.ID,
// 			IMAGE:              v.IMAGE,
// 			NAME:               v.NAME,
// 			MEMBERS:            v.MEMBERS,
// 			LOCATION_AND_DATES: delete,
// 			CREATION_DATE:      v.CREATION_DATE,
// 			FIRST_ALBUM:        v.FIRST_ALBUM,
// 			RELATIONS:          v.RELATIONS,
// 		})
// 	}

// 	return res, nil
// }
