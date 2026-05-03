package mqtt

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Subscription struct {
	Topic   string
	Qos     byte
	Handler mqtt.MessageHandler
}

type MqttClient struct {
	config MqttConfig
	client mqtt.Client
}

type MqttConfig struct {
	Broker   string
	ClientId string
	Username string
	Password string
	CAPath   string
	CertPath string
	KeyPath  string
}

func NewMqttClient(config MqttConfig) *MqttClient {
	return &MqttClient{
		config: config,
	}
}

func (m *MqttClient) Connect(defaultMessageHandler mqtt.MessageHandler, subscriptions []Subscription) error {
	//tlsconfig, err := m.NewTLSConfig(m.config.CAPath)
	// if err != nil {
	// 	return err
	// }

	log.Print("Connecting to mqtt broker...")

	// Configure MQTT client options
	opts := mqtt.NewClientOptions()
	opts.AddBroker(m.config.Broker)
	// opts.SetClientID(m.config.ClientId).SetTLSConfig(tlsconfig)

	if m.config.Username != "" {
		opts.SetUsername(m.config.Username)
	}
	if m.config.Password != "" {
		opts.SetPassword(m.config.Password)
	}

	opts.SetAutoReconnect(true)
	if defaultMessageHandler != nil {
		opts.SetDefaultPublishHandler(defaultMessageHandler)
	}
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler

	// Create and connect the client
	m.client = mqtt.NewClient(opts)

	if token := m.client.Connect(); token.WaitTimeout(5 * time.Second) {
		if token.Error() != nil {
			log.Printf("Failed to connect to MQTT broker: %v", token.Error())
			return token.Error()
		} else {
			log.Println("Successfully connected to MQTT broker")
			for _, v := range subscriptions {
				if token := m.client.Subscribe(v.Topic, v.Qos, v.Handler); token.WaitTimeout(5 * time.Second) {
					if token.Error() != nil {
						log.Printf("Failed to subscribe to topic %s: %v", v.Topic, token.Error())
						return token.Error()
					} else {
						log.Printf("Successfully subscribed to topic %s", v.Topic)
					}
				}
			}
		}
	}

	return nil
}

// Probably not going to bother using certificates for the demo, but this is a working example from older project
// if it is gonna be used, please update it to use env vars for cert paths and such
func (m *MqttClient) NewTLSConfig(cafile string) (*tls.Config, error) {
	// Import trusted certificates from CAfile.pem.
	// Alternatively, manually add CA certificates to
	// default openssl CA bundle.
	certpool := x509.NewCertPool()
	pemCerts, err := os.ReadFile(cafile)
	if err == nil {
		certpool.AppendCertsFromPEM(pemCerts)
	}

	// Import client certificate/key pair
	cert, err := tls.LoadX509KeyPair("../certs/client.crt", "../certs/client.key")
	if err != nil {
		return nil, err
	}

	// Just to print out the client certificate..
	cert.Leaf, err = x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		return nil, err
	}

	// Create tls.Config with desired tls properties
	return &tls.Config{
		// RootCAs = certs used to verify server cert.
		RootCAs: certpool,
		// ClientAuth = whether to request cert from server.
		// Since the server is set up for SSL, this happens
		// anyways.
		ClientAuth: tls.NoClientCert,
		// ClientCAs = certs used to validate client cert.
		ClientCAs: nil,
		// InsecureSkipVerify = verify that cert contents
		// match server. IP matches what is in cert etc.
		InsecureSkipVerify: true,
		// Certificates = list of certs client sends to server.
		Certificates: []tls.Certificate{cert},
	}, nil
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	log.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	log.Printf("Connect lost: %v", err)
}
