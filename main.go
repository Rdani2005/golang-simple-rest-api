package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type task struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Content string `json:"content"`
}

type allTasks []task

var tasks = allTasks{
	{
		ID:      1,
		Name:    "Task one",
		Content: "Some content",
	},
}

func getTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func createTask(w http.ResponseWriter, r *http.Request) {
	var newTask task
	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Insert a valid task")
	}
	json.Unmarshal(reqBody, &newTask)

	newTask.ID = len(tasks) + 1
	tasks = append(tasks, newTask)
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(newTask)
}

func indexRoute(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to my API on GoLang!")
}

func getTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID, error := strconv.Atoi(vars["id"])
	if error != nil {
		fmt.Fprintf(w, "Invalid ID")
		return
	}

	for _, task := range tasks {
		if task.ID == taskID {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(task)
		}
	}
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID, error := strconv.Atoi(vars["id"])
	if error != nil {
		fmt.Fprintf(w, "Invalid ID")
		return
	}

	for i, task := range tasks {
		if task.ID == taskID {
			tasks = append(tasks[:i], tasks[i+1:]...)
			fmt.Fprintf(w, "The task with ID %v has been remove successfully", taskID)
		}
	}
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", indexRoute)
	router.HandleFunc("/tasks", getTasks).Methods("GET")
	router.HandleFunc("/tasks", createTask).Methods("POST")
	router.HandleFunc("/tasks/{id}", getTask).Methods("GET")
	log.Fatal(http.ListenAndServe(":3000", router))

}
