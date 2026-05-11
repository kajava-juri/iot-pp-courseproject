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
		deviceNameStr := ""
		if len(parts) > 0 {
			deviceNameStr = parts[0]
		}
		alertTopic := deviceNameStr + "/alert/fall"

		// Lookup device by business device name (registered devices only)
		device, err := services.Device.GetByDeviceName(deviceNameStr)
		if err != nil {
			log.Printf("Unknown device %s: %v", deviceNameStr, err)
			return
		}

		// Create event and attach available device relations.
		event := &models.Event{
			Type:      models.EventTypeFall,
			DeviceID:  device.ID,
			PatientID: device.PatientID,
			Patient:   device.Patient,
			RoomID:    device.RoomID,
			Room:      device.Room,
		}

		if err := services.Event.Create(event); err != nil {
			log.Printf("Failed to create fall event: %v", err)
			wsHub.BroadcastToTopic(msg.Payload(), topic)
			return
		}

		alert := &models.Alert{
			EventID:   &event.ID,
			PatientID: device.PatientID,
			Patient:   device.Patient,
			Severity:  "high",
			Message:   fmt.Sprintf("Fall detected: %s", payload),
		}
		if err := services.Alert.Create(alert); err != nil {
			log.Printf("Failed to create fall alert: %v", err)
			wsHub.BroadcastToTopic(msg.Payload(), topic)
			return
		}

		if alert.PatientID != nil {
			status, err := services.Alert.GetPatientStatus(*alert.PatientID)
			if err != nil {
				log.Printf("Failed to derive patient %d status from active alerts: %v", *alert.PatientID, err)
			} else if patient, changed, err := services.Patient.UpdateStatus(*alert.PatientID, status); err != nil {
				log.Printf("Failed to update patient %d status: %v", *alert.PatientID, err)
			} else if changed {
				websockets.BroadcastJSONToTopic(wsHub, "/patient/update", patient)
			}
		}

		alert.Event = event
		alertJSON, err := json.Marshal(alert)
		if err != nil {
			log.Printf("Failed to marshal fall alert: %v", err)
			wsHub.BroadcastToTopic(msg.Payload(), topic)
			return
		}

		// reproadcast event
		wsHub.BroadcastToTopic(msg.Payload(), topic)
		// broadcast alert
		wsHub.BroadcastToTopic(alertJSON, alertTopic)
		client.Publish(alertTopic, 1, false, alertJSON)
	}
}

func MotionEventMessageHandler(wsHub *websockets.WsHub) mqtt.MessageHandler {
	return func(client mqtt.Client, msg mqtt.Message) {
		topic := msg.Topic()
		payload := string(msg.Payload())
		log.Printf("Motion event => %s", payload)

		parts := strings.SplitN(strings.TrimPrefix(topic, "/"), "/", 2)
		deviceNameStr := ""
		if len(parts) > 0 {
			deviceNameStr = parts[0]
		}
		// alertTopic := deviceNameStr + "/alert/motion"

		// Lookup device by business device name (registered devices only)
		device, err := services.Device.GetByDeviceName(deviceNameStr)
		if err != nil {
			log.Printf("Unknown device %s: %v", deviceNameStr, err)
			return
		}

		// Create event and attach available device relations.
		event := &models.Event{
			Type:      models.EventTypeMotionDetected,
			DeviceID:  device.ID,
			PatientID: device.PatientID,
			Patient:   device.Patient,
			RoomID:    device.RoomID,
			Room:      device.Room,
		}

		if err := services.Event.Create(event); err != nil {
			log.Printf("Failed to create motion event: %v", err)
			wsHub.BroadcastToTopic(msg.Payload(), topic)
			return
		}

		// motion alone is not enough to trigger an alert, but we broadcast the event to any interested clients
		eventJSON, err := json.Marshal(event)
		if err != nil {
			log.Printf("Failed to marshal motion event: %v", err)
			wsHub.BroadcastToTopic(msg.Payload(), topic)
			return
		}

		// if no motion, payload = 0, do not broadcast to frontend since it is not an actionable event and would only add noise
		if payload == "0" {
			return
		}
		wsHub.BroadcastToTopic(eventJSON, topic)
	}
}

func VibrationEventMessageHandler(wsHub *websockets.WsHub) mqtt.MessageHandler {
	return func(client mqtt.Client, msg mqtt.Message) {
		topic := msg.Topic()
		payload := string(msg.Payload())
		log.Printf("Vibration event => %s", payload)

		parts := strings.SplitN(strings.TrimPrefix(topic, "/"), "/", 2)
		deviceNameStr := ""
		if len(parts) > 0 {
			deviceNameStr = parts[0]
		}
		alertTopic := deviceNameStr + "/alert/vibration"

		// Lookup device by business device name (registered devices only)
		device, err := services.Device.GetByDeviceName(deviceNameStr)
		if err != nil {
			log.Printf("Unknown device %s: %v", deviceNameStr, err)
			return
		}

		// Create event and attach available device relations.
		event := &models.Event{
			Type:      "vibration",
			DeviceID:  device.ID,
			PatientID: device.PatientID,
			Patient:   device.Patient,
			RoomID:    device.RoomID,
			Room:      device.Room,
		}

		if err := services.Event.Create(event); err != nil {
			log.Printf("Failed to create vibration event: %v", err)
			wsHub.BroadcastToTopic(msg.Payload(), topic)
			return
		}

		alert := &models.Alert{
			EventID:   &event.ID,
			PatientID: device.PatientID,
			Patient:   device.Patient,
			Severity:  "medium",
			Message:   fmt.Sprintf("Vibration detected: %s", payload),
		}
		if err := services.Alert.Create(alert); err != nil {
			log.Printf("Failed to create vibration alert: %v", err)
			wsHub.BroadcastToTopic(msg.Payload(), topic)
			return
		}

		if alert.PatientID != nil {
			status, err := services.Alert.GetPatientStatus(*alert.PatientID)
			if err != nil {
				log.Printf("Failed to derive patient %d status from active alerts: %v", *alert.PatientID, err)
			} else if patient, changed, err := services.Patient.UpdateStatus(*alert.PatientID, status); err != nil {
				log.Printf("Failed to update patient %d status: %v", *alert.PatientID, err)
			} else if changed {
				websockets.BroadcastJSONToTopic(wsHub, "/patient/update", patient)
			}
		}

		alert.Event = event
		alertJSON, err := json.Marshal(alert)
		if err != nil {
			log.Printf("Failed to marshal vibration alert: %v", err)
			wsHub.BroadcastToTopic(msg.Payload(), topic)
			return
		}

		// reproadcast event
		wsHub.BroadcastToTopic(msg.Payload(), topic)
		// broadcast alert
		wsHub.BroadcastToTopic(alertJSON, alertTopic)
		client.Publish(alertTopic, 1, false, alertJSON)
	}
}

// EncoderMessageHandler returns a mqtt.MessageHandler that logs encoder payloads.
func EncoderMessageHandler() mqtt.MessageHandler {
	return func(client mqtt.Client, msg mqtt.Message) {
		log.Printf("Encoder %s => %s", msg.Topic(), string(msg.Payload()))
	}
}
