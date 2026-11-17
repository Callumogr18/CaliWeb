package server

import "net/http"

// NewRouter builds the site's routes. Endpoints from docs/api.md get added
// here as their handlers are written.
func NewRouter() *http.ServeMux {
	router := http.NewServeMux()

	// User endpoints
	router.HandleFunc("GET /{$}", home)
	router.HandleFunc("GET /projects", projects)
	router.HandleFunc("GET /thoughts", feed)

	// Admin endpoints (TODO)
	router.HandleFunc("POST /thoughts/new", createPostRouter)

	// 404
	router.HandleFunc("/", notFound)

	return router
}
