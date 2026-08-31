package main

// import (
// 	"encoding/json" //used to transfer go structs into JSON strings
// 	"fmt"           //formatting, used for
// 	"net/http"      //Built-in web library , It handles port bindings,TCP connections and route requests
// )

// type StatusResponse struct { //structure of our JSON response
// 	Status  string `json:"status"` //``are called as struct tags
// 	Message string `json:"message"`
// }

// func homeHandler(w http.ResponseWriter, r *http.Request) { //Handler for that root "/"
// 	fmt.Fprintln(w, "Welcome to my first Go web server!") //plain text response and Fprintln refers to print into a file like object
// }

// func statusHandler(w http.ResponseWriter, r *http.Request) { //Handles the api/status end point
// 	w.Header().Set("Content-Type", "application/json") //response header set to JSON

// 	response := StatusResponse{
// 		Status:  "OK",
// 		Message: "Server is running smoothly",
// 	}

// 	json.NewEncoder(w).Encode(response) //encodes our response struct into JSON and write it to the response writer 'w'
// 	//creates a JSON encoder
// }

// func main() {
// 	http.HandleFunc("/", homeHandler) //Registers handlers for routes

// 	http.HandleFunc("/api/status", statusHandler)

// 	fmt.Println("Server is starting on port 8080...")
// 	fmt.Println("open http://localhost:8080/ in your browser")

// 	err := http.ListenAndServe(":8080", nil) //starts the server
// 	if err != nil {
// 		fmt.Println("Failed to start the server", err)
// 	}

// }
