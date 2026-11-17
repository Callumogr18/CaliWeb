package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"mysite/src/store"
)

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("encode response: %v", err)
	}
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

// home serves GET /.
func home(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"page": "home"})
}

// projects serves GET /projects.
func projects(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"page": "projects"})
}

func createPostRouter(w http.ResponseWriter, r *http.Request) {
	title := strings.TrimSpace(r.FormValue("title"))
	body := strings.TrimSpace(r.FormValue("body"))

	if title == "" || body == "" {
		writeError(w, http.StatusBadRequest, "Form values cannot be null")
		return
	}

	newPost, err := store.CreatePost(title, body)
	if err != nil {
		log.Printf("create post: %v", err)
		writeError(w, http.StatusInternalServerError, "Could not create post")
		return
	}
	log.Printf("New Post: %+v", newPost)

	writeJSON(w, http.StatusCreated, newPost)
}

// feed serves /thoughts with the stored posts, newest first.
func feed(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, store.ListPosts())
}

func notFound(w http.ResponseWriter, r *http.Request) {
	writeError(w, http.StatusNotFound, "not found")
}
