package routes

import (
	"backend/database/services"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type RoomHandler struct {
}

func (h *RoomHandler) ListRooms(w http.ResponseWriter, r *http.Request) {
	rooms, err := services.Room.List()
	if err != nil {
		log.Printf("Failed to list rooms: %v", err)
		http.Error(w, "Failed to list rooms", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rooms)
}

func RoomRoutes() chi.Router {
	r := chi.NewRouter()
	roomHandler := RoomHandler{}

	r.Get("/", roomHandler.ListRooms)
	// r.Post("/", roomHandler.CreateRoom)
	// r.Get("/{id}", roomHandler.GetRoom)
	// r.Put("/{id}", roomHandler.UpdateRoom)
	// r.Delete("/{id}", roomHandler.DeleteRoom)
	return r
}
