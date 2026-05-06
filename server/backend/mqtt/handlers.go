package handlers

import (
	"backend/database/models"
	"backend/database/services"
	"backend/pkg/websockets"
	"log"
	"regexp"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// ImuMessageHandler returns a mqtt.MessageHandler that logs IMU topic payloads.
// DEPRECATED: This handler is no longer used since the IMU node now publishes fall events instead of raw sensor data.
func ImuMessageHandler() mqtt.MessageHandler {
	return func(client mqtt.Client, msg mqtt.Message) {
		r, _ := regexp.Compile("(accel|gyro)/(x|y|z)")
		if r.MatchString(msg.Topic()) {
			matches := r.FindStringSubmatch(msg.Topic())
			accelType := matches[1] // acc or gyro
			accelAxis := matches[2] // x, y, or z
			switch accelType {
			case "accel":
				log.Printf("IMU linear acceleration %s => %s", accelAxis, string(msg.Payload()))
			case "gyro":
				log.Printf("IMU angular velocity %s => %s", accelAxis, string(msg.Payload()))
			}
		}
	}
}

func FallEventMessageHandler(wsHub *websockets.WsHub) mqtt.MessageHandler {
	event := &models.Event{
		Type: models.EventTypeFall,
	}
	if err := services.Event.Create(event); err != nil {
		log.Printf("Failed to create fall event: %v", err)
	}

	return func(client mqtt.Client, msg mqtt.Message) {
		topic := msg.Topic()
		payload := string(msg.Payload())
		log.Printf("Fall event => %s", payload)
		wsHub.BroadcastToTopic(msg.Payload(), topic)
	}
}

// EncoderMessageHandler returns a mqtt.MessageHandler that logs encoder payloads.
func EncoderMessageHandler() mqtt.MessageHandler {
	return func(client mqtt.Client, msg mqtt.Message) {
		log.Printf("Encoder %s => %s", msg.Topic(), string(msg.Payload()))
	}
}
