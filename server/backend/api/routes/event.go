package routes

import (
	"backend/database/services"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type EventHandler struct {
}

func (h *EventHandler) ListEvents(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("page_size")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 50
	}

	events, err := services.Event.List(page, pageSize)
	if err != nil {
		log.Printf("Failed to list events: %v", err)
		http.Error(w, "Failed to list events", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(events); err != nil {
		log.Printf("Failed to encode events to JSON: %v", err)
	}
}

func EventRoutes() chi.Router {
	r := chi.NewRouter()
	eventHandler := EventHandler{}

	r.Get("/", eventHandler.ListEvents)
	// r.Post("/", eventHandler.CreateEvent)
	// r.Get("/{id}", eventHandler.GetEvent)
	// r.Put("/{id}", eventHandler.UpdateEvent)
	// r.Delete("/{id}", eventHandler.DeleteEvent)
	return r
}
