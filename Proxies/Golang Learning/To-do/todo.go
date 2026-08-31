package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type Task struct { //task reperesents our todo item
	ID        int    `json:"id"` //backticks tell go format the fields in lowercase
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

var (
	tasks = []Task{} //holds our todo list in memory

	mu sync.Mutex //prevents different HTTP threads from corrupting the list

	idCounter = 1 //generates unique IDs for new tasks

)

func homeHandler(w http.ResponseWriter, r *http.Request) { //handler for "/" home check
	fmt.Fprintln(w, "Welcome to the Todo App!")
}

func getTaskHandler(w http.ResponseWriter, r *http.Request) { //Handler for GET /tasks
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json") //tells the client the response will be in Json format

	mu.Lock()         //locking database before reading
	defer mu.Unlock() //unlocking database after reading

	json.NewEncoder(w).Encode(tasks) //encode slice of tasks into JSON
}

func createTaskHandler(w http.ResponseWriter, r *http.Request) { //Handler for POST /tasks/create
	if r.Method != http.MethodPost {
		http.Error(w, "Method not Allowed", http.StatusMethodNotAllowed)
		return
	}
	var input struct { //A temp struct to decode the incoming JSON payload
		Title string `json:"title"`
	}
	err := json.NewDecoder(r.Body).Decode(&input) //decodes json string and maps it to input variable
	if err != nil || input.Title == "" {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	mu.Lock()
	newTask := Task{
		ID:        idCounter,
		Title:     input.Title,
		Completed: false, //new task is always false by default
	}
	idCounter++
	tasks = append(tasks, newTask)
	mu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTask)
}

func completeTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	var input struct {
		ID int `json:"id"`
	}
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil || input.ID <= 0 {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	mu.Lock()
	defer mu.Unlock()
	for i := 0; i < len(tasks); i++ {
		if tasks[i].ID == input.ID {
			tasks[i].Completed = true
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(tasks[i])
			return
		}
	}
	http.Error(w, "Task Not Found", http.StatusNotFound)
}
func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/tasks", getTaskHandler)
	http.HandleFunc("/tasks/create", createTaskHandler)
	http.HandleFunc("/tasks/complete", completeTaskHandler)
	fmt.Println("Todo server starting on port 8080....")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Server failed to start:", err)
	}
}
