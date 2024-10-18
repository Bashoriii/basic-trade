package handlers

import (
	"basic-trade/entity"
	"basic-trade/helpers"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func (s *Service) GetAllProduct(w http.ResponseWriter, r *http.Request) {
	var products []entity.Product

	query := `SELECT id, uuid, name, image_url, admin_id FROM products`
	rows, err := s.db.Query(query)
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var product entity.Product
		if err := rows.Scan(&product.ID, &product.UUID, &product.Name, &product.ImageUrl, &product.AdminID); err != nil {
			http.Error(w, "Failed to scan product", http.StatusInternalServerError)
			return
		}
		products = append(products, product)
	}

	if err = rows.Err(); err != nil {
		http.Error(w, "Error while reading products", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func (s *Service) GetProductByUuid(w http.ResponseWriter, r *http.Request) {
	productUUID := r.FormValue("uuid")
	if productUUID == "" {
		http.Error(w, "Product ID is required", http.StatusBadRequest)
		return
	}

	var product entity.Product
	query := `SELECT id, uuid, name, image_url, admin_id FROM products WHERE uuid = $1`
	err := s.db.QueryRow(query, productUUID).Scan(&product.ID, &product.UUID, &product.Name, &product.ImageUrl, &product.AdminID)
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(product); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

func (s *Service) CreateProduct(w http.ResponseWriter, r *http.Request) {
	// Ensure the request is authenticated
	userData := r.Context().Value("userData").(map[string]interface{})
	adminID := int(userData["id"].(float64))

	// Get the form values
	productName := r.FormValue("name")
	file, fileHeader, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Error getting form file: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	if fileHeader.Size > 5<<20 {
		http.Error(w, "Max file 5MB", http.StatusBadRequest)
		return
	}

	fileType := fileHeader.Header.Get("Content-Type")
	validTypes := map[string]bool{
		"image/jpeg":    true,
		"image/jpg":     true,
		"image/png":     true,
		"image/svg+xml": true,
	}

	if !validTypes[fileType] {
		http.Error(w, "Invalid file format. Only JPG, JPEG, PNG, and SVG are allowed.", http.StatusBadRequest)
		return
	}

	// Initialize Cloudinary
	cld, err := helpers.InitCloudinary()
	if err != nil {
		http.Error(w, "Error initializing Cloudinary: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Upload the file to Cloudinary
	imageURL, err := helpers.UploadFile(cld, *fileHeader, fileHeader.Filename)
	if err != nil {
		http.Error(w, "Error uploading file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	product := entity.Product{
		Name:     productName,
		ImageUrl: imageURL,
		AdminID:  adminID,
	}

	if err := product.ValidateProduct(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Insert the product into the database
	query := `INSERT INTO products (name, image_url, admin_id) VALUES ($1, $2, $3) RETURNING id`
	err = s.db.QueryRow(query, product.Name, product.ImageUrl, product.AdminID).Scan(&product.ID)
	if err != nil {
		http.Error(w, "Unable to create product: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with success message
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Product created!"))
}

func (s *Service) UpdateProduct(w http.ResponseWriter, r *http.Request) {

	// Extract the product ID from the form-data
	productID := r.FormValue("id")
	if productID == "" {
		http.Error(w, "Product ID is required", http.StatusBadRequest)
		return
	}

	// Extract the product name from the form-data
	productName := r.FormValue("name")
	if productName == "" {
		http.Error(w, "Product name is required", http.StatusBadRequest)
		return
	}

	// Convert productID from string to int
	id, err := strconv.Atoi(productID)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	query := `UPDATE products SET name = $1 WHERE id = $2`
	_, err = s.db.Exec(query, productName, id)
	if err != nil {
		http.Error(w, "Unable to update product", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	// w.Write([]byte("Product updated!"))
	w.Write([]byte(fmt.Sprintf("Product with ID %d updated!", id)))
}

func (s *Service) UpdateProductByUuid(w http.ResponseWriter, r *http.Request) {

	// Extract the product ID from the form-data
	productUUID := r.FormValue("uuid")
	if productUUID == "" {
		http.Error(w, "Product UUID is required", http.StatusBadRequest)
		return
	}

	// Extract the product name from the form-data
	productName := r.FormValue("name")
	if productName == "" {
		http.Error(w, "Product name is required", http.StatusBadRequest)
		return
	}

	query := `UPDATE products SET name = $1 WHERE uuid = $2`
	_, err := s.db.Exec(query, productName, productUUID)
	if err != nil {
		http.Error(w, "Unable to update product", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("Product updated by UUID!")))
}

func (s *Service) DeleteProductById(w http.ResponseWriter, r *http.Request) {

	productID := r.FormValue("id")
	if productID == "" {
		http.Error(w, "Product ID is required", http.StatusBadRequest)
		return
	}

	query := `DELETE FROM products WHERE id = $1`
	_, err := s.db.Exec(query, productID)
	if err != nil {
		http.Error(w, "Unable to delete product", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Product deleted by ID!")))
}

func (s *Service) DeleteProductByUuid(w http.ResponseWriter, r *http.Request) {

	productUUID := r.FormValue("uuid")
	if productUUID == "" {
		http.Error(w, "Product UUID is required", http.StatusBadRequest)
		return
	}

	query := `DELETE FROM products WHERE uuid = $1`
	_, err := s.db.Exec(query, productUUID)
	if err != nil {
		http.Error(w, "Unable to delete product", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Product deleted by UUID!")))
}
