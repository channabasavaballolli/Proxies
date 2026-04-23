package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {
	// The port our proxy will listen on
	port := ":9000"

	// Define a handler that acts as a proxy
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("--> Forwarding request: %s %s\n", r.Method, r.URL.String())

		// 1. Create a new request to the destination
		// We use the same Method, URL, and Body as the original request
		proxyReq, err := http.NewRequest(r.Method, r.URL.String(), r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			return
		}

		// 2. Copy the headers from the original request to the proxy request
		for header, values := range r.Header {
			for _, value := range values {
				proxyReq.Header.Add(header, value)
			}
		}

		// 3. Send the request using a default client
		client := &http.Client{}
		resp, err := client.Do(proxyReq)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()

		// 4. Copy the response headers back to the original client
		for header, values := range resp.Header {
			for _, value := range values {
				w.Header().Add(header, value)
			}
		}

		// 5. Set the status code and copy the response body
		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
	})

	fmt.Printf("Forward Proxy started on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

// curl.exe -x http://localhost:9000 http://example.com
