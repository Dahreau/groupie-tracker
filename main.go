package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

type Artist struct {
	Id           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
}
type Location struct {
	Id        int      `json:"id"`
	Locations []string `json:"locations"`
}
type Coordinates struct {
	Location  string `json:"location"`
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}

type LocationsResponse struct {
	Index []Location `json:"index"`
}
type Date struct {
	Id    int      `json:"id"`
	Dates []string `json:"dates"`
}
type DateResponse struct {
	Index []Date `json:"index"`
}
type RelationsResponse struct {
	Index []Relation `json:"index"`
}
type Relation struct {
	Id             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

var (
	Artists    []Artist
	MaxMembers []int
	Locations  []Location
	Dates      []Date
	Relations  []Relation
	Coords     []Coordinates
	mu         sync.Mutex
)
var apiURLs = []string{
	"artists",
	"locations",
	"dates",
	"relation",
}

func main() {

	go FetchData()

	fmt.Println("Starting web server...")

	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	mux.HandleFunc("/", HomeHandler)
	mux.HandleFunc("/artist/", ArtistHandler)
	mux.HandleFunc("/api/artists", DataArtistsHandler)
	mux.HandleFunc("/api/dates", DataDatesHandler)
	mux.HandleFunc("/api/locations", DataLocationsHandler)
	mux.HandleFunc("/api/relations", DataRelationsHandler)

	server := &http.Server{
		Addr:              ":8080",           //adresse du server (le port choisi est à titre d'exemple)
		Handler:           mux,               // listes des handlers
		ReadHeaderTimeout: 10 * time.Second,  // temps autorisé pour lire les headers
		WriteTimeout:      10 * time.Second,  // temps maximum d'écriture de la réponse
		IdleTimeout:       120 * time.Second, // temps maximum entre deux rêquetes
		ReadTimeout:       10 * time.Second,
		MaxHeaderBytes:    1 << 20, // 1 MB // maxinmum de bytes que le serveur va lire
	}
	log.Printf("http://localhost%v", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
