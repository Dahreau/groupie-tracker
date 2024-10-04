package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

func renderTemplate(w http.ResponseWriter, tmpl string, data any) {
	t, err := template.ParseFiles("./templates/" + tmpl + ".page.tmpl")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.Execute(w, data)
}

func HomeHandler(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		renderTemplate(w, "404", nil)
	} else {
		if req.Method == "GET" {
			search := req.FormValue("search")
			searchLocation := req.FormValue("searchLocation")
			req.ParseForm()
			minDate := req.FormValue("minDate")
			maxDate := req.FormValue("maxDate")
			minDateFirstAlbum := req.FormValue("minDateFirstAlbum")
			maxDateFirstAlbum := req.FormValue("maxDateFirstAlbum")
			numberOfMembers := req.Form["numberOfMember"]
			intNumberOfMembers := make([]int, len(numberOfMembers))
			artists := Artists
			for i := 0; i < len(numberOfMembers); i++ {
				value, _ := strconv.Atoi(numberOfMembers[i])
				intNumberOfMembers = append(intNumberOfMembers, value)
			}
			if len(intNumberOfMembers) != 0 {
				artists = findArtistsByNumberOfMembers(intNumberOfMembers)
			}
			if minDate != "" && maxDate != "" {
				artists = findArtistsByMinAndMaxDate(artists, minDate, maxDate)
			}
			if minDateFirstAlbum != "" && maxDateFirstAlbum != "" {
				artists = findArtistsByMinAndMaxDateFirstAlbum(artists, minDateFirstAlbum, maxDateFirstAlbum)
			}
			searchType := req.FormValue("type")
			year := req.FormValue("year")
			date := req.FormValue("date")
			artistsName := fetchArtistsName(artists)
			locationsAvailable := findLocationsAvailable()
			creationDates := fetchCreationDates(artists)
			minCreationDate, maxCreationDate := findMinAnMaxCreationDate()
			minFirstAlbumDate, maxFirstAlbumDate := findMinAnMaxFirstAlbumDate()
			var resultArtists []Artist
			if search != "" || ((searchType == "date" || searchType == "firstAlbum") && date != "") || (searchType == "location" && searchLocation != "") || (searchType == "creationDate" && year != "") {
				switch {
				case searchType == "artist":
					resultArtists = findArtists(search, artists)
				case searchType == "date":
					resultArtists = findArtistsByDate(date, artists)
				case searchType == "firstAlbum":
					fmt.Println(date)
					resultArtists = findArtistsByFirstAlbum(date, artists)
				case searchType == "location":
					resultArtists = findArtistsByLocation(searchLocation, artists)
				case searchType == "creationDate":
					resultArtists = findArtistByCreationDate(year, artists)
				}
			} else {
				resultArtists = artists
			}
			type Data struct {
				ResultArtists      []Artist
				ArtistsName        []string
				LocationsAvailable []string
				CreationDates      []string
				MaxMembers         []int
				MinCreationDate    int
				MaxCreationDate    int
				MinDate            string
				MaxDate            string
				MinFirstAlbumDate  int
				MaxFirstAlbumDate  int
				MinDateFirstAlbum  string
				MaxDateFirstAlbum  string
			}
			data := Data{
				resultArtists,
				artistsName,
				locationsAvailable,
				creationDates,
				MaxMembers,
				minCreationDate,
				maxCreationDate,
				minDate,
				maxDate,
				minFirstAlbumDate,
				maxFirstAlbumDate,
				minDateFirstAlbum,
				maxDateFirstAlbum,
			}
			renderTemplate(w, "index", data)
		}
	}
}

func ArtistHandler(w http.ResponseWriter, req *http.Request) {

	path := req.URL.Path
	parts := strings.Split(path, "/")
	var id string
	if len(parts) >= 3 {
		id = parts[2]
	} else {
		// Handle the case where the ID is not present
		http.Error(w, "Invalid artist ID", http.StatusBadRequest)
	}
	artist := findArtistById(id)
	if artist.Id == 0 {
		// Handle the case where the artist is not found
		w.WriteHeader(http.StatusNotFound)
		renderTemplate(w, "404", nil)
		return
	}
	relations := findRelationsById(id)
	relationsFormmatted := make(map[string][]string, len(relations.DatesLocations))
	for i, relation := range relations.DatesLocations {
		tempLoc := strings.ReplaceAll(i, "_", " ")
		tempLoc = strings.ReplaceAll(tempLoc, "-", " - ")
		tempLoc = Capitalize(tempLoc)

		for _, date := range relation {
			tempDate := strings.ReplaceAll(date, "-", "/")
			relationsFormmatted[tempLoc] = append(relationsFormmatted[tempLoc], tempDate)
		}

	}
	relations.DatesLocations = relationsFormmatted
	coordArtist := []Coordinates{}
	artistLocations := []string{}
	for _, location := range Locations {
		if location.Id == artist.Id {
			artistLocations = location.Locations
		}
	}
	for _, artistLocation := range artistLocations {
		for _, coord := range Coords {
			if artistLocation == coord.Location {
				coordArtist = append(coordArtist, coord)
			}
		}
	}

	type Data struct {
		Artist      Artist
		Relations   Relation
		Coordinates []Coordinates
	}
	data := Data{
		Artist:      artist,
		Relations:   relations,
		Coordinates: coordArtist,
	}
	renderTemplate(w, "artist", data)
}

func DataArtistsHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	artists := Artists
	json.NewEncoder(w).Encode(artists)
	w.Header().Set("Content-Type", "application/json")
}
func DataDatesHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	dates := Dates
	json.NewEncoder(w).Encode(dates)
	w.Header().Set("Content-Type", "application/json")
}
func DataLocationsHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	locations := Locations
	json.NewEncoder(w).Encode(locations)
	w.Header().Set("Content-Type", "application/json")
}
func DataRelationsHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	relations := Relations
	json.NewEncoder(w).Encode(relations)
	w.Header().Set("Content-Type", "application/json")
}
func Capitalize(s string) string {
	sRune := []rune(s)
	result := []rune{}
	for i := 0; i < len(sRune); i++ {
		if i == 0 {
			if sRune[i] >= 'a' && sRune[i] <= 'z' {
				result = append(result, sRune[i]-32)
			} else {
				result = append(result, sRune[i])
			}
		} else {
			if (result[i-1] >= 'A' && result[i-1] <= 'Z') || (result[i-1] >= 'a' && result[i-1] <= 'z') || (result[i-1] >= '0' && result[i-1] <= '9') {
				if sRune[i] >= 'a' && sRune[i] <= 'z' {
					result = append(result, sRune[i])
				} else if sRune[i] >= 'A' && sRune[i] <= 'Z' {
					result = append(result, sRune[i]+32)
				} else {
					result = append(result, sRune[i])
				}
			} else {
				if sRune[i] >= 'a' && sRune[i] <= 'z' {
					result = append(result, sRune[i]-32)
				} else {
					result = append(result, sRune[i])
				}
			}
		}
	}
	return string(result)
}
