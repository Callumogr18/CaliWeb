package main

import (
	"fmt"
	"net/http"
	"os"

	server "mysite/src/server"
	"mysite/src/store"
)

func main() {
	store.Migrate()
	router := server.NewRouter()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Port for local dev
	}

	fmt.Println("Serving on port 127.0.0.1:" + port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		fmt.Println("Error running server:", err)
	}
}
