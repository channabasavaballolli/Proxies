package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func main() {
	// 1. Define where the "real" server is
	backendURL, _ := url.Parse("http://localhost:9001")

	// 2. Create the Reverse Proxy
	proxy := httputil.NewSingleHostReverseProxy(backendURL)

	// 3. (Optional) Customizing the proxy
	// Let's modify the request before it reaches the backend
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		
		// Add the 'X-Forwarded-For' header so the backend knows the client's IP
		req.Header.Set("X-Forwarded-For", req.RemoteAddr)
		
		fmt.Printf("Proxying request for: %s\n", req.URL.Path)
	}

	// 4. Start the proxy server
	port := ":9002"
	fmt.Printf("Reverse Proxy started on http://localhost%s (forwarding to %s)\n", port, backendURL)
	http.ListenAndServe(port, proxy)
}
