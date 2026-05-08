package main

import (
	"backend/api/routes"
	postgres "backend/database"
	"backend/database/services"
	handlers "backend/mqtt"
	mqttClient "backend/pkg/mqtt"
	"backend/pkg/utils"
	"backend/pkg/websockets"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type CommandArguments struct {
	CleanDb bool
}

var CmdArgs = CommandArguments{}

func main() {
	utils.ParseFlags()
	// Load environment variables from .env file
	utils.LoadEnv()

	// Initialize database connection
	if err := postgres.InitDb(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	log.Println("Database initialized")

	wsHub := websockets.NewWsHub()
	go wsHub.Run()

	mqttConfig := mqttClient.MqttConfig{
		Broker:   utils.GetEnv("MQTT_BROKER", "mqtt://193.40.245.72:1883"),
		ClientId: utils.GetEnv("MQTT_CLIENT_ID", "home-security-backend"),
		Username: utils.GetEnv("MQTT_USERNAME", "test"),
		Password: utils.GetEnv("MQTT_PASSWORD", "test"),
	}

	log.Printf("MQTT Config: %+v", mqttConfig)

	imuTopicPrefix := utils.GetEnv("WEARABLE_IMU_TOPIC_PREFIX", "ESP16")
	fallTopic := utils.GetEnv("FALL_EVENT_TOPIC", "event/fall")

	subs := []mqttClient.Subscription{
		{
			Topic:   imuTopicPrefix + "/" + fallTopic + "/#",
			Qos:     1,
			Handler: handlers.FallEventMessageHandler(wsHub),
		},
	}

	mc := mqttClient.NewMqttClient(mqttConfig)
	// TODO: default message handler?
	if err := mc.Connect(nil, subs); err != nil {
		log.Printf("mqtt connect error: %v", err)
	}

	log.Println("MQTT client connected and subscribed to topics")

	// 1. Initialize the Chi router
	r := chi.NewRouter()

	// 2. Add Middleware
	// Logger was recommended
	r.Use(middleware.Logger)
	// Add CORS middleware to allow requests from the frontend
	// Credit: https://github.com/go-chi/cors
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	websockets.RegisterRoutes(r, wsHub)

	// 3. Create a Dummy API (GET Endpoint)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	})

	r.Get("/event/{id}", func(w http.ResponseWriter, r *http.Request) {
		idParam := chi.URLParam(r, "id")
		eventID, err := strconv.ParseUint(idParam, 10, 64)
		if err != nil {
			http.Error(w, "Invalid event ID", http.StatusBadRequest)
			return
		}

		event, err := services.Event.GetByID(uint(eventID))
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(event); err != nil {
			log.Printf("Failed to encode event to JSON: %v", err)
		}
	})

	r.Get("/event/all", func(w http.ResponseWriter, r *http.Request) {
		events, err := services.Event.GetAll()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(events); err != nil {
			log.Printf("Failed to encode events to JSON: %v", err)
		}
	})

	r.Mount("/alert", routes.AlertRoutes(wsHub))
	r.Mount("/patient", routes.PatientRoutes())
	r.Mount("/room", routes.RoomRoutes())
	r.Mount("/event", routes.EventRoutes())

	// 4. Run the Server
	// Start server in goroutine
	go func() {
		if err := http.ListenAndServe(":"+utils.GetEnv("WEB_PORT", "8080"), r); err != nil {
			log.Printf("API server error: %v", err)
		}
	}()

	// keep the main function running to allow MQTT client to receive messages
	select {}
}
