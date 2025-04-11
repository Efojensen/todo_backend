package main

import (
	"net/http"
)

type Todo struct {
	Id       uint8  `json:"id"`
	Task     string `json:"task"`
	Complete bool   `json:"complete"`
}

func getAllTodos (w http.ResponseWriter, r *http.Request) {

}

func createTodo (w http.ResponseWriter, r *http.Request) {

}

func updateTodo (w http.ResponseWriter, r *http.Request) {

}

func deleteTodo (w http.ResponseWriter, r *http.Request) {

}

func main() {
	// Read Route
	http.HandleFunc("/", getAllTodos)

	// Create Route
	http.HandleFunc("/create", createTodo)

	// Update Route
	http.HandleFunc("/{id}", updateTodo)

	// Delete Route
	http.HandleFunc("{id}", deleteTodo)

	http.ListenAndServe(":4000", nil)
}
