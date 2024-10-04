package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func FetchData() {
	mu.Lock()
	defer mu.Unlock()
	fmt.Println("lalalalalalalalalal")
	// Fetch data from external APIs
	for _, name := range apiURLs {
		switch name {
		case "artists":
			Artists = fetchArtists()
		case "locations":
			Locations = fetchLocations()
		case "dates":
			Dates = fetchDates()
		case "relation":
			Relations = fetchRelations()
		}
	}
	WriteCoord()
	fmt.Println("Data fetched successfully")
}
func WriteCoord() {
	locations := Locations
	type ApiResponse struct {
		Lat string `json:"lat"`
		Lon string `json:"lon"`
	}
	file, err := os.OpenFile("static/coord.json", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(1)
	}
	defer file.Close()
	var existingCoords []Coordinates
	err = json.NewDecoder(file).Decode(&existingCoords)
	if err != nil {
		fmt.Println("Error decoding JSON1:", err)
		os.Exit(1)
	}
	var coords []ApiResponse
	for i := 0; i < len(locations); i++ {
		artistLoc := locations[i]
		for _, loc := range artistLoc.Locations {
			isExisting := false
			for _, existingCoord := range existingCoords {
				if loc == existingCoord.Location {
					isExisting = true
					break
				}
			}
			if !isExisting {

				resp, err := http.Get("https://nominatim.openstreetmap.org/search?q=" + strings.ReplaceAll(loc, "-", ",") + "&format=json&limit=1")
				if err != nil {
					fmt.Println("Error:", err)
					os.Exit(1)
				}
				defer resp.Body.Close()
				err = json.NewDecoder(resp.Body).Decode(&coords)
				if err != nil {
					fmt.Println("Error decoding JSON2:", err)
					os.Exit(1)
				}
				if len(coords) > 0 && coords[0].Lat != "" && coords[0].Lon != "" {
					coordinates := Coordinates{Location: loc, Latitude: coords[0].Lat, Longitude: coords[0].Lon}
					existingCoords = append(existingCoords, coordinates)
				}
				if err != nil {
					fmt.Println("Error writing file:", err)
					os.Exit(1)
				}
			}
		}
	}
	err = file.Truncate(0)
	if err != nil {
		fmt.Println("Error truncating file:", err)
		os.Exit(1)
	}
	file.Seek(0, 0)
	Coords = existingCoords
	// Write all coordinates to the file
	err = json.NewEncoder(file).Encode(existingCoords)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		os.Exit(1)
	}
}
func fetchArtists() []Artist {
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error: status code", resp.StatusCode)
		os.Exit(1)
	}
	var artists []Artist
	err = json.NewDecoder(resp.Body).Decode(&artists)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		os.Exit(1)
	}
	maxMembers := 0
	for _, artist := range artists {
		if len(artist.Members) > maxMembers {
			maxMembers = len(artist.Members)
		}
	}
	for i := 1; i <= maxMembers; i++ {
		MaxMembers = append(MaxMembers, i)
	}
	return artists
}

func fetchLocations() []Location {
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/locations")
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error: status code", resp.StatusCode)
		os.Exit(1)
	}
	var locationResponse LocationsResponse
	err = json.NewDecoder(resp.Body).Decode(&locationResponse)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		os.Exit(1)
	}
	var locations []Location
	locations = append(locations, locationResponse.Index...)
	return locations
}
func fetchDates() []Date {
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/dates")
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error: status code", resp.StatusCode)
		os.Exit(1)
	}
	var dateResponse DateResponse
	err = json.NewDecoder(resp.Body).Decode(&dateResponse)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		os.Exit(1)
	}
	var dates []Date
	dates = append(dates, dateResponse.Index...)
	return dates
}
func fetchRelations() []Relation {
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/relation")
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error: status code", resp.StatusCode)
		os.Exit(1)
	}
	var relationsResponse RelationsResponse
	err = json.NewDecoder(resp.Body).Decode(&relationsResponse)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		os.Exit(1)
	}
	var relations []Relation
	relations = append(relations, relationsResponse.Index...)
	return relations
}
func findArtistById(id string) Artist {
	artists := Artists
	var result Artist
	for i := range artists {
		if strconv.Itoa(artists[i].Id) == id {
			result = artists[i]
		}
	}
	return result

}
func findRelationsById(id string) Relation {
	relations := Relations
	var result Relation
	for i := range relations {
		if strconv.Itoa(relations[i].Id) == id {
			result = relations[i]
		}
	}
	return result
}
func findArtists(search string, artists []Artist) []Artist {
	var result []Artist

	searchFormatted := strings.ReplaceAll(search, " (Artist/Band)", "")
	searchFormatted = strings.ReplaceAll(searchFormatted, " (Member)", "")

	for _, artist := range artists {
		if strings.HasPrefix(strings.ToLower(artist.Name), strings.ToLower(searchFormatted)) {
			result = append(result, artist)
		} else {
			for _, member := range artist.Members {
				if strings.HasPrefix(strings.ToLower(member), strings.ToLower(searchFormatted)) {
					result = append(result, artist)
					break
				}
			}
		}

	}
	return result
}

func findArtistsByDate(date string, artists []Artist) []Artist {
	var result []Artist
	var resultIds []int
	datesArtist := fetchDates()
	dateTbl := strings.Split(date, "-")
	dateFormatted := ""
	for i := len(dateTbl) - 1; i >= 0; i-- {
		dateFormatted += dateTbl[i] + "-"
	}
	dateFormatted = dateFormatted[:len(dateFormatted)-1]

	for _, dates := range datesArtist {
		for _, date := range dates.Dates {

			if strings.Contains(date, dateFormatted) {
				for _, artist := range artists {
					if artist.Id == dates.Id && !contains(resultIds, artist.Id) {
						resultIds = append(resultIds, artist.Id)
						result = append(result, artist)
						break
					}
				}
			}

		}
	}
	return result
}
func findArtistsByMinAndMaxDate(artists []Artist, min string, max string) []Artist {
	minDate, _ := strconv.Atoi(min)
	maxDate, _ := strconv.Atoi(max)
	var result []Artist
	for _, artist := range artists {
		if artist.CreationDate >= minDate && artist.CreationDate <= maxDate {
			result = append(result, artist)
		}
	}
	return result
}
func findArtistsByMinAndMaxDateFirstAlbum(artists []Artist, min string, max string) []Artist {
	minDate, _ := strconv.Atoi(min)
	maxDate, _ := strconv.Atoi(max)
	var result []Artist
	for _, artist := range artists {
		yearFirstAlbum, _ := strconv.Atoi(artist.FirstAlbum[len(artist.FirstAlbum)-4:])
		if yearFirstAlbum >= minDate && yearFirstAlbum <= maxDate {
			result = append(result, artist)
		}
	}
	return result
}
func findArtistsByFirstAlbum(date string, artists []Artist) []Artist {
	var result []Artist
	fmt.Println(date)
	dateTbl := strings.Split(date, "-")
	dateFormatted := ""
	for i := len(dateTbl) - 1; i >= 0; i-- {
		dateFormatted += dateTbl[i] + "-"
	}
	fmt.Println(dateFormatted)
	dateFormatted = dateFormatted[:len(dateFormatted)-1]
	for _, artist := range artists {
		if strings.Contains(artist.FirstAlbum, dateFormatted) {
			result = append(result, artist)
		}
	}
	return result
}

func findArtistByCreationDate(year string, artists []Artist) []Artist {
	var result []Artist
	fmt.Println("Year : ", year)
	yearFormatted := year[0:4]
	fmt.Println("yearFormatted:", yearFormatted)
	for _, artist := range artists {
		if strings.Contains(strconv.Itoa(artist.CreationDate), yearFormatted) {
			result = append(result, artist)
		}
	}
	return result
}

func findArtistsByNumberOfMembers(numbers []int) []Artist {
	artists := Artists
	var result []Artist
	for _, artist := range artists {
		for i := range numbers {
			if len(artist.Members) == numbers[i] {
				result = append(result, artist)
				break
			}
		}
	}
	return result
}
func fetchCreationDates(artists []Artist) []string {
	var result []string
	for _, artist := range artists {
		year := strconv.Itoa(artist.CreationDate) + " - " + artist.Name
		result = append(result, year)
	}
	return result
}

func findArtistsByLocation(search string, artists []Artist) []Artist {
	var result []Artist
	var resultIds []int
	searchFormatted := strings.ToLower(search)
	searchFormatted = strings.TrimSpace(searchFormatted)
	searchFormatted = strings.ReplaceAll(searchFormatted, " - ", "-")
	searchFormatted = strings.ReplaceAll(searchFormatted, " ", "_")
	locationsArtist := fetchLocations()
	for _, locations := range locationsArtist {
		for _, location := range locations.Locations {
			if strings.Contains(location, searchFormatted) {
				for _, artist := range artists {
					if artist.Id == locations.Id && !contains(resultIds, artist.Id) {
						resultIds = append(resultIds, artist.Id)
						result = append(result, artist)
						break
					}
				}
			}

		}
	}
	return result
}

func fetchArtistsName(artists []Artist) []string {
	var result []string
	for _, artist := range artists {
		artistName := artist.Name + " (Artist/Band)"
		if !containsMembers(result, artistName) {
			result = append(result, artistName)
		}
		for _, member := range artist.Members {
			memberName := member + " (Member)"
			if !containsMembers(result, memberName) {
				result = append(result, memberName)
			}
		}
	}
	return result
}

func findLocationsAvailable() []string {
	locations := Locations
	var result []string
	for _, locationArtist := range locations {
		for _, location := range locationArtist.Locations {
			locationFormatted := strings.ReplaceAll(location, "_", " ")
			locationFormatted = strings.ReplaceAll(locationFormatted, "-", " - ")
			locationFormatted = Capitalize(locationFormatted)
			if !containsLocations(result, locationFormatted) {
				result = append(result, locationFormatted)
			}
		}
	}
	return result
}
func findMinAnMaxCreationDate() (int, int) {
	artists := Artists
	minDate := 9999999999999
	maxDate := 0
	for _, artist := range artists {
		if artist.CreationDate < minDate {
			minDate = artist.CreationDate
		}
		if artist.CreationDate > maxDate {
			maxDate = artist.CreationDate
		}
	}
	return minDate, maxDate

}
func findMinAnMaxFirstAlbumDate() (int, int) {
	artists := Artists
	minDate := 9999999999999
	maxDate := 0
	for _, artist := range artists {
		yearFirstAlbum, _ := strconv.Atoi(artist.FirstAlbum[len(artist.FirstAlbum)-4:])
		if yearFirstAlbum < minDate {
			minDate = yearFirstAlbum
		}
		if yearFirstAlbum > maxDate {
			maxDate = yearFirstAlbum
		}
	}
	return minDate, maxDate

}
func containsMembers(result []string, memberToAdd string) bool {
	for _, member := range result {
		if member == memberToAdd {
			return true
		}
	}
	return false
}
func containsLocations(result []string, locationToAdd string) bool {
	for _, location := range result {
		if location == locationToAdd {
			return true
		}
	}
	return false
}
func contains(slice []int, number int) bool {
	for _, v := range slice {
		if v == number {
			return true
		}
	}
	return false
}
