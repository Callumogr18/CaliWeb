package server

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"mysite/src/store"
)

// HTTP handlers. Every GET route answers with an HTML page; the one JSON
// endpoint left is the admin POST at the bottom.

// serverError logs err and serves the plain 500 page.
func serverError(w http.ResponseWriter, err error) {
	log.Printf("server error: %v", err)
	render(w, http.StatusInternalServerError, "error", newPage("Error", ""))
}

// showHome serves GET /.
func showHome(w http.ResponseWriter, r *http.Request) {
	posts, err := store.ListPosts()
	if err != nil {
		serverError(w, err)
		return
	}
	data := HomePage{Page: newPage("Home", "home")}
	if len(posts) > 0 {
		data.Latest = &posts[0]
	}
	render(w, http.StatusOK, "home", data)
}

// showProjects serves GET /projects.
func showProjects(w http.ResponseWriter, r *http.Request) {
	render(w, http.StatusOK, "projects", newPage("Projects", "projects"))
}

// showThoughts serves GET /thoughts with the stored posts, newest first.
func showThoughts(w http.ResponseWriter, r *http.Request) {
	posts, err := store.ListPosts()
	if err != nil {
		serverError(w, err)
		return
	}
	render(w, http.StatusOK, "thoughts", ThoughtsPage{Page: newPage("Thoughts", "thoughts"), Posts: posts})
}

// showPost serves GET /thoughts/{slug}.
func showPost(w http.ResponseWriter, r *http.Request) {
	post, err := store.GetPostBySlug(r.PathValue("slug"))
	switch {
	case errors.Is(err, store.ErrNotFound):
		showNotFound(w, r)
		return
	case err != nil:
		serverError(w, err)
		return
	}
	render(w, http.StatusOK, "post", PostPage{
		Page:       newPage(post.Title, "thoughts"),
		Post:       post,
		Paragraphs: paragraphs(post.Body),
	})
}

// showNotFound is the catch-all "/" handler and the unknown-slug response.
func showNotFound(w http.ResponseWriter, r *http.Request) {
	render(w, http.StatusNotFound, "notfound", newPage("Not found", ""))
}

// createPostRouter serves POST /thoughts/new. It stays JSON because there is
// no UI for writing yet.
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
