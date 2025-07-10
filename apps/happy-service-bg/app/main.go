package main

import (
	"log"
	"net/http"
	"os"
)

const version = "v1"

func main() {
	port := ":8080"

	// Serve static files
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/static", home)
	http.HandleFunc("/ping", pingHandler) // Add the ping handler

	loggedMux := loggingMiddleware(http.DefaultServeMux)

	log.Printf("Starting server on %s...\n", port)
	if err := http.ListenAndServe(port, loggedMux); err != nil {
		log.Fatalf("Error: %v\n", err)
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		next.ServeHTTP(w, r)
	})
}

func home(w http.ResponseWriter, r *http.Request) {
	file := "static/index_" + version + ".html"

	if _, err := os.Stat(file); os.IsNotExist(err) {
		http.Error(w, "404 - Page Not Found", http.StatusNotFound)
		return
	}

	http.ServeFile(w, r, file)
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong"))
}
