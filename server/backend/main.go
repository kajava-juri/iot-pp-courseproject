package main

import (
	postgres "backend/database"
	handlers "backend/mqtt"
	mqttClient "backend/pkg/mqtt"
	"backend/pkg/utils"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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

	mqttConfig := mqttClient.MqttConfig{
		Broker:   utils.GetEnv("MQTT_BROKER", "mqtt://193.40.245.72:1883"),
		ClientId: utils.GetEnv("MQTT_CLIENT_ID", "home-security-backend"),
		Username: utils.GetEnv("MQTT_USERNAME", "test"),
		Password: utils.GetEnv("MQTT_PASSWORD", "test"),
	}

	log.Printf("MQTT Config: %+v", mqttConfig)

	// publishReading("sensor/imu/accel/x", a.acceleration.x);
	// publishReading("sensor/imu/accel/y", a.acceleration.y);
	// publishReading("sensor/imu/accel/z", a.acceleration.z);
	// publishReading("sensor/imu/gyro/x", g.gyro.x);
	// publishReading("sensor/imu/gyro/y", g.gyro.y);
	// publishReading("sensor/imu/gyro/z", g.gyro.z);

	imuTopicPrefix := utils.GetEnv("WEARABLE_IMU_TOPIC_PREFIX", "ESP16")

	subs := []mqttClient.Subscription{
		{
			Topic:   imuTopicPrefix + "/imu/#",
			Qos:     1,
			Handler: handlers.ImuMessageHandler(),
		},
		{
			Topic:   imuTopicPrefix + "/enc",
			Qos:     1,
			Handler: handlers.EncoderMessageHandler(),
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

	// 2. Add Middleware (optional but recommended)
	r.Use(middleware.Logger)

	// 3. Create a Dummy API (GET Endpoint)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	})

	// 4. Run the Server
	http.ListenAndServe(":8080", r)

	// keep the main function running to allow MQTT client to receive messages
	select {}
}
