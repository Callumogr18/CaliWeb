package main

import (
	"fmt"
	"net/http"

	server "mysite/src/server"
	"mysite/src/store"
)

func main() {
	store.Migrate()
	router := server.NewRouter()

	fmt.Println("Serving on http://localhost:8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		fmt.Println("Error running server:", err)
	}
}
