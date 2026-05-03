package handlers

import (
	"log"
	"regexp"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// ImuMessageHandler returns a mqtt.MessageHandler that logs IMU topic payloads.
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

// EncoderMessageHandler returns a mqtt.MessageHandler that logs encoder payloads.
func EncoderMessageHandler() mqtt.MessageHandler {
	return func(client mqtt.Client, msg mqtt.Message) {
		log.Printf("Encoder %s => %s", msg.Topic(), string(msg.Payload()))
	}
}
