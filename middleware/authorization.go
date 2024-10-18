package middleware

import (
	"basic-trade/database"
	"basic-trade/entity"
	"database/sql"
	"net/http"
	"strconv"
)

func ProductAuthorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		productUUID := r.FormValue("uuid")
		if productUUID == "" {
			http.Error(w, "Product UUID is required", http.StatusBadRequest)
			return
		}

		// Get admin ID from the request context (set by Authentication middleware)
		userData := r.Context().Value("userData").(map[string]interface{})
		adminID := int(userData["id"].(float64))

		// Query the product from the database to check the admin who created it
		db := database.DBConnection()
		defer db.Close()

		var product entity.Product
		query := `SELECT admin_id FROM products WHERE uuid = $1`
		err := db.QueryRow(query, productUUID).Scan(&product.AdminID)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Product not found", http.StatusNotFound)
				return
			}
			http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Check if the admin ID matches
		if product.AdminID != adminID {
			http.Error(w, "Unauthorized: You do not have access to modify this product", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func VariantAuthorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Validate only Admin who created the product that can add variant to it by ID
		productIdStr := r.FormValue("product_id")
		if productIdStr == "" {
			http.Error(w, "Product ID is required (Dari authorization)", http.StatusBadRequest)
			return
		}

		productID, err := strconv.Atoi(productIdStr)
		if err != nil {
			http.Error(w, "Invalid Product ID", http.StatusBadRequest)
			return
		}

		// Get admin ID from the request context (set by Authentication middleware)
		userData := r.Context().Value("userData").(map[string]interface{})
		adminID := int(userData["id"].(float64))

		// Query the product from the database to check the admin who created it
		db := database.DBConnection()
		defer db.Close()

		var product entity.Product
		query := `SELECT admin_id FROM products WHERE id = $1`
		err = db.QueryRow(query, productID).Scan(&product.AdminID)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Product not found", http.StatusNotFound)
				return
			}
			http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		if product.AdminID != adminID {
			http.Error(w, "Unauthorized: You do not have access to modify this product", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func UpdateVariantAuthorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		variantUUID := r.FormValue("uuid")
		if variantUUID == "" {
			http.Error(w, "Product UUID is required", http.StatusBadRequest)
			return
		}

		// Get admin ID from the request context (set by Authentication middleware)
		userData := r.Context().Value("userData").(map[string]interface{})
		adminID := int(userData["id"].(float64))

		// Query the product associated with the variant to check the admin who created it
		db := database.DBConnection()
		defer db.Close()

		var productID int
		query := `SELECT p.admin_id FROM variants v 
				  JOIN products p ON v.product_id = p.id 
				  WHERE v.uuid = $1`

		err := db.QueryRow(query, variantUUID).Scan(&productID)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Variant not found", http.StatusNotFound)
				return
			}
			http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Check if the admin ID matches the creator of the product
		if productID != adminID {
			http.Error(w, "Unauthorized: You do not have access to modify this variant", http.StatusUnauthorized)
			return
		}

		// Allow the request to proceed
		next.ServeHTTP(w, r)
	})
}

func DeleteVariantAuthorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Extract UUID from URL path (assuming chi router)
		variantUUID := r.FormValue("uuid")
		if variantUUID == "" {
			http.Error(w, "Variant UUID is required", http.StatusBadRequest)
			return
		}

		// Get admin ID from the request context (set by Authentication middleware)
		userData := r.Context().Value("userData").(map[string]interface{})
		adminID := int(userData["id"].(float64))

		// Query the product associated with the variant to check the admin who created it
		db := database.DBConnection()
		defer db.Close()

		var productAdminID int
		query := `
            SELECT p.admin_id 
            FROM variants v 
            JOIN products p ON v.product_id = p.id 
            WHERE v.uuid = $1`

		err := db.QueryRow(query, variantUUID).Scan(&productAdminID)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Variant not found", http.StatusNotFound)
				return
			}
			http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Check if the requesting admin matches the product's creator
		if productAdminID != adminID {
			http.Error(w, "Unauthorized: You do not have access to delete this variant", http.StatusUnauthorized)
			return
		}

		// Proceed to the next handler
		next.ServeHTTP(w, r)
	})
}
