package handlers

import (
	"basic-trade/entity"
	"fmt"
	"net/http"
	"strconv"
)

func (s *Service) CreateVariant(w http.ResponseWriter, r *http.Request) {

	variantName := r.FormValue("variant_name")
	variantQty := r.FormValue("qty")
	productIdStr := r.FormValue("product_id")

	quantity, err := strconv.Atoi(variantQty)
	if err != nil {
		http.Error(w, "Invalid quantity value", http.StatusBadRequest)
		return
	}

	productID, err := strconv.Atoi(productIdStr)
	if err != nil {
		http.Error(w, "Invalid quantity value", http.StatusBadRequest)
		return
	}

	variant := entity.Variant{
		VariantName: variantName,
		Quantity:    quantity,
		ProductID:   productID,
	}

	if err := variant.ValidateVariant(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query := `INSERT INTO variants (variant_name, quantity, product_id) VALUES ($1, $2, $3) RETURNING id`
	err = s.db.QueryRow(query, variant.VariantName, variant.Quantity, variant.ProductID).Scan(&variant.ID)
	if err != nil {
		http.Error(w, "Unable to create variant", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Variant created!"))
}

func (s *Service) UpdateVariant(w http.ResponseWriter, r *http.Request) {

	variantUUID := r.FormValue("uuid")
	if variantUUID == "" {
		http.Error(w, "Variant UUID is required", http.StatusBadRequest)
		return
	}

	variantName := r.FormValue("variant_name")
	if variantName == "" {
		http.Error(w, "Variant name is requid", http.StatusBadRequest)
		return
	}

	query := `UPDATE variants SET variant_name = $1 WHERE uuid = $2`
	_, err := s.db.Exec(query, variantName, variantUUID)
	if err != nil {
		http.Error(w, "Unable to update variant", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Variant updated by UUID!"))
}

func (s *Service) DeleteVariantByUuid(w http.ResponseWriter, r *http.Request) {

	variantUUID := r.FormValue("uuid")
	if variantUUID == "" {
		http.Error(w, "Variant UUID is required", http.StatusBadRequest)
		return
	}

	query := `DELETE FROM variants WHERE uuid = $1`
	_, err := s.db.Exec(query, variantUUID)
	if err != nil {
		http.Error(w, "Unable to delete variant", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Product deleted by UUID!")))
}
