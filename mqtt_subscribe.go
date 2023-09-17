package main

import (
    "fmt"
    "os"
    "os/signal"
    "strings"
    "github.com/eclipse/paho.mqtt.golang"
)

func onMessageReceived(client mqtt.Client, message mqtt.Message) {
    fmt.Printf("Received message on topic '%s': %s\n", message.Topic(), message.Payload())
}

