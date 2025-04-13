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

func createTodo(w http.ResponseWriter, r *http.Request) {
	var todo Todo
	err := json.NewDecoder(r.Body).Decode(&todo)

	if err != nil {
		http.Error(w, "Error reading from the request body", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	res, err := collection.InsertOne(r.Context(), todo)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	var new Todo
	err = collection.FindOne(r.Context(), bson.M{"_id": res.InsertedID}).Decode(&new)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err = json.NewEncoder(w).Encode(new); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func updateTodo(w http.ResponseWriter, r *http.Request) {

}

func deleteTodo(w http.ResponseWriter, r *http.Request) {

}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading csv file. Err: ", err)
	}

	MONGODB_URI := os.Getenv("MONGODB_URI")
	clientOptions := options.Client().ApplyURI(MONGODB_URI)
	client, err := mongo.Connect(context.Background(), clientOptions)

	if err != nil {
		log.Println("Err: ", err)
	}

	defer client.Disconnect(context.Background())

	if err := client.Ping(context.Background(), nil); err != nil {
		log.Println("Database unable to be pinged")
	}

	fmt.Println("Connection to mongo client successful")

	collection = client.Database("Career-Atlas").Collection("todos_v2")

	// READ Route
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var todos []Todo
		cursor, err := collection.Find(context.Background(), bson.M{})

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		defer cursor.Close(context.Background())

		var todo Todo
		for cursor.Next(context.Background()) {
			if err := cursor.Decode(&todo); err != nil {
				http.Error(w, "Invalid structure\nActual err: ", http.StatusNotFound)
			}
			todos = append(todos, todo)
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(todos); err != nil {
			log.Printf("Failed to encode todos: %v", err)
		}
	})

	// Create Route
	http.HandleFunc("/create", createTodo)

	// Update Route
	http.HandleFunc("/{id}", updateTodo)

	// Delete Route
	http.HandleFunc("/delete/{id}", deleteTodo)

	http.ListenAndServe(":4000", nil)
}
