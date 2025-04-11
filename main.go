package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Todo struct {
	Task     string `json:"task"`
	Complete bool   `json:"complete"`
}

var collection *mongo.Collection

func createTodo (w http.ResponseWriter, r *http.Request) {
	var todo Todo
	err := json.NewDecoder(r.Body).Decode(&todo)

	if err != nil {
		http.Error(w, "Error reading from the request body", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	newTodo, err := collection.InsertOne(context.Background(), todo)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
	}

	json.NewEncoder(w).Encode(newTodo)
	w.WriteHeader(http.StatusCreated)
}

func updateTodo (w http.ResponseWriter, r *http.Request) {

}

func deleteTodo (w http.ResponseWriter, r *http.Request) {

}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading csv file. Err: ", err)
	}

	MONGODB_URI := os.Getenv("MONGODB_URI")
	clientOptions := options.Client().ApplyURI(MONGODB_URI)
	client, err := mongo.Connect(context.Background(), clientOptions)

	if err != nil {
		log.Fatal("Err: ", err)
	}

	if err := client.Ping(context.Background(), nil); err != nil {
		log.Fatal("Database unable to be pinged")
	}

	fmt.Println("Connection to mongo client successful")

	collection := client.Database("Career-Atlas").Collection("todos_v2")

	// Read Route
	http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
		var todos []Todo
		cursor, err := collection.Find(context.Background(), bson.M{})

		if err != nil {
			log.Fatal("Err: ", err)
		}

		defer cursor.Close(context.Background())

		var todo Todo
		for cursor.Next(context.Background()) {
			if err := cursor.Decode(&todo); err != nil {
				log.Fatal("Invalid structure\nActual err: ", err)
			}
			todos = append(todos, todo)
		}

		jsonData, err := json.Marshal(todos)
		if err != nil {
			log.Fatal("Error unmarshaling data")
		}

		json.NewEncoder(w).Encode(jsonData)
	})

	// Create Route
	http.HandleFunc("/create", createTodo)

	// Update Route
	http.HandleFunc("/{id}", updateTodo)

	// Delete Route
	http.HandleFunc("/delete/{id}", deleteTodo)

	http.ListenAndServe(":4000", nil)
}
