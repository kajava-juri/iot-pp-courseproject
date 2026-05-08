package routes

import (
	"backend/database/models"
	"backend/database/services"
	"backend/pkg/websockets"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

type AlertHandler struct {
	hub *websockets.WsHub
}

func NewAlertHandler(hub *websockets.WsHub) *AlertHandler {
	return &AlertHandler{hub: hub}
}

func BroadcastAlertUpdate(hub *websockets.WsHub, alert *models.Alert) {
	if hub == nil {
		log.Printf("Websocket hub is nil, cannot broadcast alert update")
		return
	}

	b, err := json.Marshal(alert)
	if err != nil {
		log.Printf("Failed to marshal alert for broadcasting: %v", err)
		return
	}

	hub.BroadcastToTopic(b, "alert/updates")
}

func (h *AlertHandler) ListUnresolvedAlerts(w http.ResponseWriter, r *http.Request) {
	alerts, err := services.Alert.GetAll(&models.Alert{Resolved: false})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	unresolvedAlerts := make([]models.Alert, 0)
	for _, alert := range alerts {
		if !alert.Resolved {
			unresolvedAlerts = append(unresolvedAlerts, alert)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(unresolvedAlerts); err != nil {
		log.Printf("Failed to encode alerts to JSON: %v", err)
	}
}

func (h *AlertHandler) AcknowledgeAlert(w http.ResponseWriter, r *http.Request) {
	alert, err := alertFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	alert.Acknowledged = true
	if err := services.Alert.Update(alert); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	updatedAlert, err := services.Alert.GetByID(alert.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(updatedAlert); err != nil {
		log.Printf("Failed to encode acknowledged alert to JSON: %v", err)
	}

	// Broadcast the updated alert to websocket clients
	BroadcastAlertUpdate(h.hub, updatedAlert)
}

// Decline alert
func (h *AlertHandler) DeclineAlert(w http.ResponseWriter, r *http.Request) {
	alert, err := alertFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	alert.Declined = true
	if err := services.Alert.Update(alert); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	updatedAlert, err := services.Alert.GetByID(alert.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(updatedAlert); err != nil {
		log.Printf("Failed to encode declined alert to JSON: %v", err)
	}

	// Broadcast the updated alert to websocket clients
	BroadcastAlertUpdate(h.hub, updatedAlert)
}

func (h *AlertHandler) ResolveAlert(w http.ResponseWriter, r *http.Request) {
	alert, err := alertFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	now := time.Now().UTC()
	alert.Resolved = true
	alert.ResolvedAt = &now
	if err := services.Alert.Update(alert); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	updatedAlert, err := services.Alert.GetByID(alert.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(updatedAlert); err != nil {
		log.Printf("Failed to encode resolved alert to JSON: %v", err)
	}

	// Broadcast the updated alert to websocket clients
	BroadcastAlertUpdate(h.hub, updatedAlert)
}

func AlertRoutes(hub *websockets.WsHub) chi.Router {
	r := chi.NewRouter()
	alertHandler := NewAlertHandler(hub)

	r.Get("/all/unresolved", alertHandler.ListUnresolvedAlerts)
	r.Post("/{id}/acknowledge", alertHandler.AcknowledgeAlert)
	r.Post("/{id}/resolve", alertHandler.ResolveAlert)
	r.Post("/{id}/decline", alertHandler.DeclineAlert)

	return r
}

func alertFromRequest(r *http.Request) (*models.Alert, error) {
	idParam := chi.URLParam(r, "id")
	alertID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return nil, err
	}

	alert, err := services.Alert.GetByID(uint(alertID))
	if err != nil {
		return nil, err
	}

	return alert, nil
}
