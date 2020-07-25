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
	"github.com/viniciusbsneto/star-wars-rest-api/framework/storage"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Planet struct
type Planet struct {
	ID      primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name    string             `json:"name,omitempty" bson:"name,omitempty"`
	Climate []string           `json:"climate,omitempty" bson:"climate,omitempty"`
	Terrain []string           `json:"terrain,omitempty" bson:"terrain,omitempty"`
	Films   int                `json:"films,omitempty" bson:"films,omitempty"`
}

// A SWAPIResponse struct to map the entire SWAPI response
type SWAPIResponse struct {
	Results []SWAPIFilms `json:"results"`
}

// A SWAPIFilms struct to map every film to
type SWAPIFilms struct {
	Films []string `json:"films"`
}

func getSWAPIPlanet(planetName string, w http.ResponseWriter, r *http.Request) (int, error) {

	response, err := http.Get("https://swapi.dev/api/planets?search=" + planetName)

	if err != nil {
		w.WriteHeader(response.StatusCode)
		w.Write([]byte(`{ "error": ` + err.Error() + ` }`))
		return 0, err
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		w.WriteHeader(response.StatusCode)
		w.Write([]byte(`{ "error": ` + err.Error() + ` }`))
		return 0, err
	}

	defer response.Body.Close()

	var swapiResponse SWAPIResponse
	err = json.Unmarshal(responseData, &swapiResponse)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{ "error": "Planet name does not exist."}`))
		return 0, err
	}

	films := len(swapiResponse.Results[0].Films)

	return films, err
}

func getPlanets(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.Header().Set("Content-Type", "application/json")
	var planets []Planet
	collection := storage.DB.Collection("planets")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		responseWriter.WriteHeader(http.StatusInternalServerError)
		responseWriter.Write([]byte(`{ "error": ` + err.Error() + ` }`))
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
		responseWriter.Write([]byte(`{ "error": ` + err.Error() + ` }`))
		return
	}
	responseWriter.WriteHeader(http.StatusOK)
	err = json.NewEncoder(responseWriter).Encode(planets)
}

func getPlanetByID(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.Header().Set("Content-Type", "application/json")
	params := mux.Vars(request)
	id, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		responseWriter.WriteHeader(http.StatusBadRequest)
		responseWriter.Write([]byte(`{ "error": ` + err.Error() + ` }`))
		return
	}
	var planet Planet
	collection := storage.DB.Collection("planets")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = collection.FindOne(ctx, bson.M{"_id": id}).Decode(&planet)
	if err != nil {
		responseWriter.WriteHeader(http.StatusBadRequest)
		responseWriter.Write([]byte(`{ "error": "Planet ID does not exist." }`))
		return
	}
	responseWriter.WriteHeader(http.StatusOK)
	json.NewEncoder(responseWriter).Encode(planet)
}

func getPlanetByName(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.Header().Set("Content-Type", "application/json")
	name := request.FormValue("name")
	var planet Planet
	collection := storage.DB.Collection("planets")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := collection.FindOne(ctx, bson.M{"name": name}).Decode(&planet)
	if err != nil {
		responseWriter.WriteHeader(http.StatusInternalServerError)
		responseWriter.Write([]byte(`{ "error": ` + err.Error() + ` }`))
		return
	}
	responseWriter.WriteHeader(http.StatusOK)
	json.NewEncoder(responseWriter).Encode(planet)
}

func createPlanet(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.Header().Set("Content-Type", "application/json")
	var planet Planet
	err := json.NewDecoder(request.Body).Decode(&planet)
	if err != nil {
		responseWriter.WriteHeader(http.StatusBadRequest)
		responseWriter.Write([]byte(`{ "error": ` + err.Error() + ` }`))
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := storage.DB.Collection("planets")
	var findOneResult bson.M
	err = collection.FindOne(ctx, bson.M{"name": planet.Name}).Decode(&findOneResult)
	if err == nil {
		responseWriter.WriteHeader(http.StatusBadRequest)
		responseWriter.Write([]byte(`{ "error": "Planet ` + planet.Name + ` already exists." }`))
		return
	}
	if err != mongo.ErrNoDocuments {
		responseWriter.WriteHeader(http.StatusBadRequest)
		responseWriter.Write([]byte(`{ "error": ` + err.Error() + ` }`))
		return
	}
	planet.Films, err = getSWAPIPlanet(planet.Name, responseWriter, request)
	if err != nil {
		return
	}
	_, err = collection.InsertOne(ctx, planet)
	if err != nil {
		responseWriter.WriteHeader(http.StatusInternalServerError)
		responseWriter.Write([]byte(`{ "error": ` + err.Error() + ` }`))
		return
	}
	responseWriter.WriteHeader(http.StatusCreated)
	json.NewEncoder(responseWriter).Encode(planet)
}

func updatePlanet(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.Header().Set("Content-Type", "application/json")
	params := mux.Vars(request)
	id, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		responseWriter.WriteHeader(http.StatusBadRequest)
		responseWriter.Write([]byte(`{ "error": ` + err.Error() + ` }`))
		return
	}
	var planet Planet
	err = json.NewDecoder(request.Body).Decode(&planet)
	if err != nil {
		responseWriter.WriteHeader(http.StatusBadRequest)
		responseWriter.Write([]byte(`{ "error": ` + err.Error() + ` }`))
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := storage.DB.Collection("planets")
	var findOneResult bson.M
	err = collection.FindOne(ctx, bson.M{"_id": id}).Decode(&findOneResult)
	if err != nil {
		responseWriter.WriteHeader(http.StatusBadRequest)
		responseWriter.Write([]byte(`{ "error": "Planet ID does not exist." }`))
		return
	}
	_, err = collection.ReplaceOne(
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
		responseWriter.Write([]byte(`{ "error": ` + err.Error() + ` }`))
		return
	}
	responseWriter.WriteHeader(http.StatusOK)
	json.NewEncoder(responseWriter).Encode(planet)
}

func deletePlanet(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.Header().Set("Content-Type", "application/json")
	params := mux.Vars(request)
	id, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		responseWriter.WriteHeader(http.StatusBadRequest)
		responseWriter.Write([]byte(`{ "error": ` + err.Error() + ` }`))
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := storage.DB.Collection("planets")
	var findOneResult bson.M
	err = collection.FindOne(ctx, bson.M{"_id": id}).Decode(&findOneResult)
	if err != nil {
		responseWriter.WriteHeader(http.StatusBadRequest)
		responseWriter.Write([]byte(`{ "error": "Planet ID does not exist." }`))
		return
	}
	_, err = collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		responseWriter.WriteHeader(http.StatusInternalServerError)
		responseWriter.Write([]byte(`{ "error": ` + err.Error() + ` }`))
		return
	}
	responseWriter.WriteHeader(http.StatusNoContent)
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
		panic(err)
	}

	router := mux.NewRouter()

	router.HandleFunc("/planets", getPlanets).Methods("GET")
	router.HandleFunc("/planets/{id}", getPlanetByID).Methods("GET")
	router.HandleFunc("/search", getPlanetByName).Methods("GET")
	router.HandleFunc("/planets", createPlanet).Methods("POST")
	router.HandleFunc("/planets/{id}", updatePlanet).Methods("PUT")
	router.HandleFunc("/planets/{id}", deletePlanet).Methods("DELETE")

	server := &http.Server{
		Handler:      router,
		Addr:         (os.Getenv("HOST") + ":" + os.Getenv("PORT")),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(server.ListenAndServe())
}
