package main

import (
	"basic-trade/database"
	"basic-trade/handlers"
	"basic-trade/middleware"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {

	db := database.DBConnection()
	defer db.Close()
	r := chi.NewRouter()

	service := handlers.NewService(db)

	r.Get("/products", service.GetAllProduct)
	r.Get("/products/{uuid}", service.GetProductByUuid)

	// Public routes
	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", service.RegisterAdmin)
		r.Post("/login", service.LoginAdmin)
	})

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(middleware.Authentication)

		r.Post("/products", service.CreateProduct)
		// r.Put("/{id}", service.UpdateProduct)
		r.With(middleware.ProductAuthorization).Put("/products/{uuid}", service.UpdateProductByUuid)
		// r.Delete("/{id}", service.DeleteProductById)
		r.With(middleware.ProductAuthorization).Delete("/products/{uuid}", service.DeleteProductByUuid)

		// Variants
		r.With(middleware.VariantAuthorization).Post("/products/variants", service.CreateVariant)

		r.With(middleware.UpdateVariantAuthorization).Put("/products/variants/{uuid}", service.UpdateVariant)

		r.With(middleware.UpdateVariantAuthorization).Delete("/products/variants/{uuid}", service.DeleteVariantByUuid)
	})

	// Check database connection
	// var currentTime string
	// err := db.QueryRow("SELECT NOW()").Scan(&currentTime)
	// if err != nil {
	// 	log.Fatalf("Failed to execute test query: %v\n", err)
	// }

	// fmt.Printf("Database connected successfully! Current time: %s\n", currentTime)
	// Check database connection

	fmt.Println("Server is starting on port 8080...")
	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}
