package main

import (
	"log"
	"net/http"
	"os"
	"time"
)

const version = "v1"

// customResponseWriter captures status code for logging
type customResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *customResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

// loggingMiddleware wraps all handlers to log method, path, status, and timing
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lrw := &customResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(lrw, r)
		duration := time.Since(start)

		log.Printf("%s %s %d %s", r.Method, r.URL.Path, lrw.statusCode, duration)
	})
}

func main() {
	port := ":8080"

	// Serve static files
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Wrap home and ping handlers
	http.Handle("/static", loggingMiddleware(http.HandlerFunc(home)))
	http.Handle("/ping", loggingMiddleware(http.HandlerFunc(pingHandler)))

	// Wrap default mux for all handlers (including static)
	loggedMux := loggingMiddleware(http.DefaultServeMux)

	log.Printf("Starting server on %s...\n", port)
	if err := http.ListenAndServe(port, loggedMux); err != nil {
		log.Fatalf("Error: %v\n", err)
	}
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
