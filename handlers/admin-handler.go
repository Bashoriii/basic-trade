package handlers

import (
	"basic-trade/entity"
	"basic-trade/helpers"
	"database/sql"
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	db *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{db: db}
}

func (s *Service) LoginAdmin(w http.ResponseWriter, r *http.Request) {

	var admin = entity.Admin{}

	admin.Email = r.FormValue("email")
	admin.Password = r.FormValue("password")

	// Retrieve id, email sama password
	query := `SELECT id, email, password FROM admins where email = $1`
	row := s.db.QueryRow(query, admin.Email)

	err := row.Scan(&admin.ID, &admin.Email, &admin.Password)
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(r.FormValue("password")))

	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	token := helpers.GenerateToken(uint(admin.ID), admin.Email)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})
}

func (s *Service) RegisterAdmin(w http.ResponseWriter, r *http.Request) {

	var admin = entity.Admin{}

	admin.Name = r.FormValue("name")
	admin.Email = r.FormValue("email")
	admin.Password = r.FormValue("password")

	if err := admin.ValidateAdmin(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query := `INSERT INTO admins (name, email, password) VALUES ($1, $2, $3) RETURNING id`
	err := s.db.QueryRow(query, admin.Name, admin.Email, admin.Password).Scan(&admin.ID)
	if err != nil {
		http.Error(w, "Unable to register", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Admin registered successfully!"))
}

func (s *Service) UpdateAdminTest(w http.ResponseWriter, r *http.Request) {

	var admin = entity.Admin{}

	if err := json.NewDecoder(r.Body).Decode(&admin); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query := `UPDATE admins SET name = $1 WHERE id = $2`
	_, err := s.db.Exec(query, admin.Name, admin.ID)
	if err != nil {
		http.Error(w, "Unable to update admin", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Admin updated!"))
}
