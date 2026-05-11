package routes

import (
	"backend/database/models"
	"backend/database/services"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type DeviceHandler struct {
}

func NewDeviceHandler() *DeviceHandler {
	return &DeviceHandler{}
}

func (h *DeviceHandler) ListDevices(w http.ResponseWriter, r *http.Request) {
	devices, err := services.Device.List()
	if err != nil {
		log.Printf("Failed to list devices: %v", err)
		http.Error(w, "Failed to list devices", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(devices)
}

func (h *DeviceHandler) UpdateDevice(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	deviceID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		http.Error(w, "Invalid device ID", http.StatusBadRequest)
		return
	}

	var payload models.Device
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Ensure ID matches URL
	payload.ID = uint(deviceID)

	// Prevent changing device name by ensuring DeviceName is empty in updates
	payload.DeviceName = ""

	if err := services.Device.Update(&payload); err != nil {
		log.Printf("Failed to update device %d: %v", payload.ID, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	updated, err := services.Device.GetByID(payload.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updated)
}

func DeviceRoutes() chi.Router {
	r := chi.NewRouter()
	handler := NewDeviceHandler()

	r.Get("/", handler.ListDevices)
	r.Put("/{id}", handler.UpdateDevice)

	return r
}
