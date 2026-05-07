package handlers

import (
	"backend/database/models"
	"backend/database/services"
	"backend/pkg/websockets"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"

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
	return func(client mqtt.Client, msg mqtt.Message) {
		topic := msg.Topic()
		payload := string(msg.Payload())
		log.Printf("Fall event => %s", payload)

		parts := strings.SplitN(strings.TrimPrefix(topic, "/"), "/", 2)
		topicPrefix := ""
		if len(parts) > 0 {
			topicPrefix = parts[0]
		}
		alertTopic := topicPrefix + "/alert/fall"

		event := &models.Event{Type: models.EventTypeFall}
		if err := services.Event.Create(event); err != nil {
			log.Printf("Failed to create fall event: %v", err)
			wsHub.BroadcastToTopic(msg.Payload(), topic)
			return
		}

		alert := &models.Alert{
			EventID:  &event.ID,
			Severity: "high",
			Message:  fmt.Sprintf("Fall detected: %s", payload),
		}
		if err := services.Alert.Create(alert); err != nil {
			log.Printf("Failed to create fall alert: %v", err)
			wsHub.BroadcastToTopic(msg.Payload(), topic)
			return
		}

		alert.Event = event
		alertJSON, err := json.Marshal(alert)
		if err != nil {
			log.Printf("Failed to marshal fall alert: %v", err)
			wsHub.BroadcastToTopic(msg.Payload(), topic)
			return
		}

		wsHub.BroadcastToTopic(msg.Payload(), topic)
		wsHub.BroadcastToTopic(alertJSON, alertTopic)
	}
}

// EncoderMessageHandler returns a mqtt.MessageHandler that logs encoder payloads.
func EncoderMessageHandler() mqtt.MessageHandler {
	return func(client mqtt.Client, msg mqtt.Message) {
		log.Printf("Encoder %s => %s", msg.Topic(), string(msg.Payload()))
	}
}
