package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-ini/ini"
)

func onMessageReceived(client mqtt.Client, message mqtt.Message) {
	// Print the message to the console with the time of the reception, and topic name:
	fmt.Printf("[%s] %s\n", time.Now().Format(time.RFC3339), message.Payload())
}

func getConfigFromIniFile() (string, string, error) {
	cfg, err := ini.Load("./config.ini")
	if err != nil {
		return "", "", err
	} else {
		if cfg.Section("DEFAULT").Key("server").String() != "" && cfg.Section("DEFAULT").Key("topic").String() != "" {
			server := cfg.Section("DEFAULT").Key("server").String()
			topic := cfg.Section("DEFAULT").Key("topic").String()
			return server, topic, nil
		} else {
			return "", "", fmt.Errorf("missing settings in config file")
		}
	}
}

func getConfigFromEnvironment() (string, string, error) {
	// Get the values of the environment variables
	if os.Getenv("Mqtt_subscribe_server") == "" || os.Getenv("Mqtt_subscribe_topics") == "" {
		return "", "", fmt.Errorf("missing environment variables Mqtt_subscribe_server or Mqtt_subscribe_topics")
	} else {
		server := os.Getenv("Mqtt_subscribe_server")
		topic := os.Getenv("Mqtt_subscribe_topics")
		return server, topic, nil
	}
}

func getConfigFromArguments() (string, string, error) {
	// Get the values of the command-line arguments
	if len(os.Args) == 3 {
		server := os.Args[1] // MQTT broker URI (e.g., "tcp://broker.hivemq.com:1883")
		topic := os.Args[2]  // MQTT topic to subscribe to
		return server, topic, nil
	} else {
		return "", "", fmt.Errorf("missing command-line arguments")
	}
}

// getServerAndTopic returns the MQTT server and topic to subscribe to.
// It first checks if the server and topic are provided as command line arguments.
// If not, it checks if they are present in a config.ini file.
// If not, it checks if they are set as environment variables.
// If none of the above options are available, it returns an error.
func getServerAndTopic() (string, string, error) {
	var server, topic string
	var err error

	switch {
	case len(os.Args) == 3:
		server, topic, err = getConfigFromArguments()
	case fileExists("config.ini"):
		server, topic, err = getConfigFromIniFile()
	case os.Getenv("Mqtt_subscribe_server") != "" && os.Getenv("Mqtt_subscribe_topics") != "":
		server, topic, err = getConfigFromEnvironment()
	default:
		err = fmt.Errorf("no server or topic found")
	}

	return server, topic, err
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// main function starts the MQTT subscriber
func main() {
	// Inside main function
	fmt.Println("Starting MQTT subscriber...")

	// Get the MQTT broker URI and topic name
	server, topic, err := getServerAndTopic()
	if err != nil {
		fmt.Printf("Error getting server and topic: %v\n", err)
		os.Exit(1)
	} else {
		fmt.Printf("Connecting to %s and subscribing to topic '%s'\n", server, topic)
	}

	// Create an MQTT client options struct
	opts := mqtt.NewClientOptions()
	opts.AddBroker(server)

	// Create an MQTT client instance
	client := mqtt.NewClient(opts)

	// Connect to the MQTT broker
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		fmt.Printf("Error connecting to MQTT broker: %v\n", token.Error())
		os.Exit(1)
	} else {
		// Print a message indicating that we're connected with the MQTT broker and display the client ID and broker URI
		fmt.Printf("Connected to MQTT broker with client ID '%s' and broker URI '%s'\n", opts.ClientID, server)
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
	fmt.Println("Listening to topic. Press Ctrl+C to disconnect.")
	<-sigCh
	fmt.Println("CTRL+C received.")

	// Unsubscribe and disconnect from the MQTT broker
	if token := client.Unsubscribe(topic); token.Wait() && token.Error() != nil {
		fmt.Printf("Error unsubscribing from topic '%s': %v\n", topic, token.Error())
	} else {
		// Print a message indicating that we're unsubscribed from the MQTT topic
		fmt.Printf("Unsubscribed from topic '%s'\n", topic)
	}

	client.Disconnect(250)

	fmt.Println("MQTT disconnected")
}
