package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

// Planet struct
type Planet struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Climate []string `json:"climate"`
	Terrain []string `json:"terrain"`
	Films   int      `json:"films"`
}

var planets []Planet

func getPlanets(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.Header().Set("Content-Type", "application/json")
	json.NewEncoder(responseWriter).Encode(planets)
}

func getPlanetByID(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.Header().Set("Content-Type", "application/json")
	params := mux.Vars(request)
	for _, item := range planets {
		if item.ID == params["id"] {
			json.NewEncoder(responseWriter).Encode(item)
			return
		}
	}
	json.NewEncoder(responseWriter).Encode(&Planet{})
}

func getPlanetByName(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.Header().Set("Content-Type", "application/json")
	name := request.FormValue("name")
	fmt.Println(name)
	for _, item := range planets {
		if item.Name == name {
			json.NewEncoder(responseWriter).Encode(item)
			return
		}
	}
	json.NewEncoder(responseWriter).Encode(&Planet{})
}

func createPlanet(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.Header().Set("Content-Type", "application/json")
	var planet Planet
	_ = json.NewDecoder(request.Body).Decode(&planet)
	planet.ID = uuid.NewV4().String()
	planets = append(planets, planet)
	json.NewEncoder(responseWriter).Encode(planet)
}

func updatePlanet(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.Header().Set("Content-Type", "application/json")
	params := mux.Vars(request)
	for index, item := range planets {
		if item.ID == params["id"] {
			planets = append(planets[:index], planets[index+1:]...)
			var planet Planet
			_ = json.NewDecoder(request.Body).Decode(&planet)
			planet.ID = params["id"]
			planets = append(planets, planet)
			json.NewEncoder(responseWriter).Encode(planet)
			return
		}
	}
	json.NewEncoder(responseWriter).Encode(planets)
}

func deletePlanet(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.Header().Set("Content-Type", "application/json")
	params := mux.Vars(request)
	for index, item := range planets {
		if item.ID == params["id"] {
			planets = append(planets[:index], planets[index+1:]...)
			break
		}
	}
	json.NewEncoder(responseWriter).Encode(planets)
}

func main() {

	router := mux.NewRouter()

	planets = append(planets, Planet{ID: "1", Name: "Tatooine", Climate: []string{"Arid"}, Terrain: []string{"Dessert"}, Films: 5})
	planets = append(planets, Planet{ID: "2", Name: "Alderaan", Climate: []string{"Temperate"}, Terrain: []string{"Grasslands", "Mountain"}, Films: 2})

	router.HandleFunc("/", func(responseWriter http.ResponseWriter, request *http.Request) {
		json.NewEncoder(responseWriter).Encode(map[string]string{"message": "Hello World"})
	}).Methods("GET")

	router.HandleFunc("/planets", getPlanets).Methods("GET")
	router.HandleFunc("/planets/{id}", getPlanetByID).Methods("GET")
	router.HandleFunc("/search", getPlanetByName).Methods("GET")
	router.HandleFunc("/planets", createPlanet).Methods("POST")
	router.HandleFunc("/planets/{id}", updatePlanet).Methods("PUT")
	router.HandleFunc("/planets/{id}", deletePlanet).Methods("DELETE")

	server := &http.Server{
		Handler:      router,
		Addr:         "localhost:3333",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(server.ListenAndServe())
}
