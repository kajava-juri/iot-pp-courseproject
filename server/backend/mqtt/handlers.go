package handlers

import (
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// ImuMessageHandler returns a mqtt.MessageHandler that logs IMU topic payloads.
func ImuMessageHandler() mqtt.MessageHandler {
	return func(client mqtt.Client, msg mqtt.Message) {
		log.Printf("IMU %s => %s", msg.Topic(), string(msg.Payload()))
	}
}

// EncoderMessageHandler returns a mqtt.MessageHandler that logs encoder payloads.
func EncoderMessageHandler() mqtt.MessageHandler {
	return func(client mqtt.Client, msg mqtt.Message) {
		log.Printf("Encoder %s => %s", msg.Topic(), string(msg.Payload()))
	}
}
