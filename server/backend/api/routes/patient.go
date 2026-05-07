package routes

import (
	"backend/database/models"
	"backend/database/services"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(patient); err != nil {
		log.Printf("Failed to encode patient to JSON: %v", err)
	}
}

func PatientRoutes() chi.Router {
	r := chi.NewRouter()
	patientHandler := PatientHandler{}

	// r.Get("/", patientHandler.ListPatients)
	r.Post("/", patientHandler.CreatePatient)
	// r.Get("/{id}", patientHandler.GetPatient)
	// r.Put("/{id}", patientHandler.UpdatePatient)
	// r.Delete("/{id}", patientHandler.DeletePatient)
	return r
}
