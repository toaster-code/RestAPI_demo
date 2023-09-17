package main

import (
	"fmt"
	"os"
	"os/signal"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func onMessageReceived(client mqtt.Client, message mqtt.Message) {
	fmt.Printf("Received message on topic '%s': %s\n", message.Topic(), message.Payload())
}

func main() {
	// Inside main function
	fmt.Println("Starting MQTT subscriber...")

	// Check for the required command-line arguments
	if len(os.Args) != 3 {
		fmt.Println("Usage: mqtt_subscribe <broker_uri> <topic>")
		os.Exit(1)
	}

	brokerURI := os.Args[1] // MQTT broker URI (e.g., "tcp://broker.hivemq.com:1883")
	topic := os.Args[2]     // MQTT topic to subscribe to

	// Print the MQTT broker URI and topic name
	fmt.Printf("Connecting to %s and subscribing to topic '%s'\n", brokerURI, topic)

	// Create an MQTT client options struct
	opts := mqtt.NewClientOptions()
	opts.AddBroker(brokerURI)

	// Create an MQTT client instance
	client := mqtt.NewClient(opts)

	// Connect to the MQTT broker
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		fmt.Printf("Error connecting to MQTT broker: %v\n", token.Error())
		os.Exit(1)
	} else {
		// Print a message indicating that we're connected with the MQTT broker and display the client ID and broker URI
		fmt.Printf("Connected to MQTT broker with client ID '%s' and broker URI '%s'\n", opts.ClientID, brokerURI)
	}

	// Subscribe to the specified MQTT topic and set the callback function
	if token := client.Subscribe(topic, 0, onMessageReceived); token.Wait() && token.Error() != nil {
		fmt.Printf("Error subscribing to topic '%s': %v\n", topic, token.Error())
		os.Exit(1)
	} else {
		// Print a message indicating that we're subscribed to the MQTT topic
		fmt.Printf("Subscribed to topic '%s'\n", topic)
	}

	// Listen for SIGINT (Ctrl+C) to gracefully exit the program
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	<-sigCh

	// Unsubscribe and disconnect from the MQTT broker
	if token := client.Unsubscribe(topic); token.Wait() && token.Error() != nil {
		fmt.Printf("Error unsubscribing from topic '%s': %v\n", topic, token.Error())
	} else {
		// Print a message indicating that we're unsubscribed from the MQTT topic
		fmt.Printf("Unsubscribed from topic '%s'\n", topic)
	}

	client.Disconnect(250)

	fmt.Println("MQTT subscriber disconnected")

}
