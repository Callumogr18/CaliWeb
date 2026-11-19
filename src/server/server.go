package server

import "net/http"

// NewRouter builds the site's routes. Endpoints from docs/api.md get added
// here as their handlers are written.
func NewRouter() *http.ServeMux {
	router := http.NewServeMux()

	// User pages (HTML)
	router.HandleFunc("GET /{$}", showHome)
	router.HandleFunc("GET /projects", showProjects)
	router.HandleFunc("GET /thoughts", showThoughts)
	router.HandleFunc("GET /thoughts/{slug}", showPost)

	// Stylesheet for the pages
	router.HandleFunc("GET /static/style.css", styleSheet)

	// Admin endpoints (TODO). Still JSON: there is no UI for writing yet.
	router.HandleFunc("POST /thoughts/new", createPostRouter)

	// 404 page
	router.HandleFunc("/", showNotFound)

	return router
}
