func main() {
    // Check for the required command-line arguments
    if len(os.Args) != 3 {
        fmt.Println("Usage: mqtt_subscribe <broker_uri> <topic>")
        os.Exit(1)
    }

    brokerURI := os.Args[1]  // MQTT broker URI (e.g., "tcp://broker.hivemq.com:1883")
    topic := os.Args[2]      // MQTT topic to subscribe to

    // Create an MQTT client options struct
    opts := mqtt.NewClientOptions()
    opts.AddBroker(brokerURI)

    // Create an MQTT client instance
    client := mqtt.NewClient(opts)

    // Connect to the MQTT broker
    if token := client.Connect(); token.Wait() && token.Error() != nil {
        fmt.Printf("Error connecting to MQTT broker: %v\n", token.Error())
        os.Exit(1)
    }

    // Subscribe to the specified MQTT topic and set the callback function
    if token := client.Subscribe(topic, 0, onMessageReceived); token.Wait() && token.Error() != nil {
        fmt.Printf("Error subscribing to topic '%s': %v\n", topic, token.Error())
        os.Exit(1)
    }

    // Listen for SIGINT (Ctrl+C) to gracefully exit the program
    sigCh := make(chan os.Signal, 1)
    signal.Notify(sigCh, os.Interrupt)
    <-sigCh

    // Unsubscribe and disconnect from the MQTT broker
    if token := client.Unsubscribe(topic); token.Wait() && token.Error() != nil {
        fmt.Printf("Error unsubscribing from topic '%s': %v\n", topic, token.Error())
    }

    client.Disconnect(250)
}
