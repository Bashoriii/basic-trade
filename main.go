package main

import (
	"basic-trade/database"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	db := database.DBConnection()
	defer db.Close()

	r := chi.NewRouter()

	fmt.Println("Server is starting on port 8080...")
	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}
