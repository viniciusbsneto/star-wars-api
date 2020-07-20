package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Planet struct
type Planet struct {
	ID      primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name    string             `json:"name" bson:"name,omitempty"`
	Climate []string           `json:"climate" bson:"climate,omitempty"`
	Terrain []string           `json:"terrain" bson:"terrain,omitempty"`
	Films   int                `json:"films" bson:"films,omitempty"`
}

// A SWAPIResponse struct to map the entire SWAPI response
type SWAPIResponse struct {
	Results []SWAPIFilms `json:"results"`
}

// A SWAPIFilms struct to map every film to
type SWAPIFilms struct {
	Films []string `json:"films"`
}

func getSWAPIPlanet(planetName string) int {

	response, err := http.Get("https://swapi.dev/api/planets?search=" + planetName)

	if err != nil {
		log.Fatalf("The HTTP request to SWAPI failed with error: %v", err)
	}

	responseData, _ := ioutil.ReadAll(response.Body)

	defer response.Body.Close()

	var swapiResponse SWAPIResponse
	json.Unmarshal(responseData, &swapiResponse)

	films := len(swapiResponse.Results[0].Films)

	return films
}

func getPlanets(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.Header().Set("Content-Type", "application/json")
	var planets []Planet
	collection := client.Database("starwars").Collection("planets")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		responseWriter.WriteHeader(http.StatusInternalServerError)
		responseWriter.Write([]byte(`{ "error": }` + err.Error() + `"}`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var planet Planet
		cursor.Decode(&planet)
		planets = append(planets, planet)
	}
	if err := cursor.Err(); err != nil {
		responseWriter.WriteHeader(http.StatusInternalServerError)
		responseWriter.Write([]byte(`{ "error": }` + err.Error() + `"}`))
		return
	}
	json.NewEncoder(responseWriter).Encode(planets)
}

func getPlanetByID(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.Header().Set("Content-Type", "application/json")
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var planet Planet
	collection := client.Database("starwars").Collection("planets")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&planet)
	if err != nil {
		responseWriter.WriteHeader(http.StatusInternalServerError)
		responseWriter.Write([]byte(`{ "error": }` + err.Error() + `"}`))
		return
	}
	json.NewEncoder(responseWriter).Encode(planet)
}

func getPlanetByName(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.Header().Set("Content-Type", "application/json")
	name := request.FormValue("name")
	var planet Planet
	collection := client.Database("starwars").Collection("planets")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err := collection.FindOne(ctx, bson.M{"name": name}).Decode(&planet)
	if err != nil {
		responseWriter.WriteHeader(http.StatusInternalServerError)
		responseWriter.Write([]byte(`{ "error": }` + err.Error() + `"}`))
		return
	}
	json.NewEncoder(responseWriter).Encode(planet)
}

func createPlanet(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.Header().Set("Content-Type", "application/json")
	var planet Planet
	json.NewDecoder(request.Body).Decode(&planet)
	planet.Films = getSWAPIPlanet(planet.Name)
	collection := client.Database("starwars").Collection("planets")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	result, _ := collection.InsertOne(ctx, planet)
	json.NewEncoder(responseWriter).Encode(result)
}

func updatePlanet(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.Header().Set("Content-Type", "application/json")
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var planet Planet
	json.NewDecoder(request.Body).Decode(&planet)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	collection := client.Database("starwars").Collection("planets")
	result, err := collection.ReplaceOne(
		ctx,
		bson.M{"_id": id},
		bson.M{
			"name":    planet.Name,
			"climate": planet.Climate,
			"terrain": planet.Terrain,
			"films":   planet.Films,
		},
	)
	if err != nil {
		responseWriter.WriteHeader(http.StatusInternalServerError)
		responseWriter.Write([]byte(`{ "error": }` + err.Error() + `"}`))
		return
	}
	json.NewEncoder(responseWriter).Encode(result)
}

func deletePlanet(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.Header().Set("Content-Type", "application/json")
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	collection := client.Database("starwars").Collection("planets")
	result, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		responseWriter.WriteHeader(http.StatusInternalServerError)
		responseWriter.Write([]byte(`{ "error": }` + err.Error() + `"}`))
		return
	}
	json.NewEncoder(responseWriter).Encode(result)
}

// ConnectDB returns Client and Context to connect to database
func ConnectDB() (*mongo.Client, context.Context) {

	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	uri := os.Getenv("MONGODB_URI")

	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("Error creating a new MongoDB client: %v", err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %v", err)
		panic(err)
	}

	return client, ctx
}

var client, ctx = ConnectDB()

func main() {

	router := mux.NewRouter()

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
