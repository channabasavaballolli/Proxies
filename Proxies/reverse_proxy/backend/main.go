package main

import (
	"flag"
	"fmt"
	"net/http"
)

// ==========================================
// 1. SERVICE LAYER (The Brain)
// ==========================================

type Greeter interface {
	GetWelcomeMessage(ip string) string
}

type GreetingService struct {
	Port string
}

func (s GreetingService) GetWelcomeMessage(ip string) string {
	return fmt.Sprintf("Hello! You have reached the Backend Server on Port %s. \n Your IP is: %s", s.Port, ip)
}

// ==========================================
// 2. HANDLER LAYER (The Face)
// ==========================================

type AppHandler struct {
	Service Greeter
}

func (h AppHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = r.RemoteAddr
	}

	message := h.Service.GetWelcomeMessage(ip)
	
	// Add this log line so we can see it in the terminal!
	fmt.Printf("--> [%s] Received request from %s\n", r.Method, ip)

	fmt.Fprintln(w, message)
}

// ==========================================
// 3. MAIN (Connecting the layers)
// ==========================================

func main() {
	// Add a flag so we can specify the port when running the app
	portPtr := flag.String("port", "9001", "Port to listen on")
	flag.Parse()
	port := ":" + *portPtr

	// Initialize our Service WITH the port name
	service := GreetingService{
		Port: *portPtr,
	}

	handler := AppHandler{
		Service: service,
	}

	fmt.Printf("Backend Server started on http://localhost%s\n", port)
	http.ListenAndServe(port, handler)
}
