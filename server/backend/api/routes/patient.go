package routes

import (
	"backend/database/models"
	"backend/database/services"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

type PatientHandler struct {
}

func (h *PatientHandler) CreatePatient(w http.ResponseWriter, r *http.Request) {
	var patient models.Patient
	err := json.NewDecoder(r.Body).Decode(&patient)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := services.Patient.Create(&patient); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			http.Error(w, "Patient with the same ID already exists", http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(patient); err != nil {
		log.Printf("Failed to encode patient to JSON: %v", err)
	}
}

func (h *PatientHandler) ListPatients(w http.ResponseWriter, r *http.Request) {
	patients, err := services.Patient.List()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(patients); err != nil {
		log.Printf("Failed to encode patients to JSON: %v", err)
	}
}

// func (h *PatientHandler) SearchPatients(w http.ResponseWriter, r *http.Request) {
// 	query := r.URL.Query().Get("q")
// 	if query == "" {
// 		http.Error(w, "Missing search query", http.StatusBadRequest)
// 		return
// 	}

func PatientRoutes() chi.Router {
	r := chi.NewRouter()
	patientHandler := PatientHandler{}

	r.Get("/", patientHandler.ListPatients)
	// r.Get("/search", patientHandler.SearchPatients)
	r.Post("/", patientHandler.CreatePatient)
	// r.Get("/{id}", patientHandler.GetPatient)
	// r.Put("/{id}", patientHandler.UpdatePatient)
	// r.Delete("/{id}", patientHandler.DeletePatient)
	return r
}
