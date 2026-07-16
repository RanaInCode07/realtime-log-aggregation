package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"realtime-log-aggregation/internal/models"
	"time"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

// Mock Log generator

var levels = []string{"INFO", "DEBUG", "WARN", "ERROR"}
var services = []string{"pg-router", "watchtower", "tokens", "terminals", "order-service", "account-service", "splitz"}
var messages = []string{
	"User authentication successful",
    "Database connection timeout retrying...",
    "Failed to process payment: insufficient funds",
    "Cache miss for product inventory metadata",
    "API gateway rate limit exceeded for client IP",
}
var messageCount = 0

func main (){
	// initializing kafka writer configuration 
	writer := &kafka.Writer{
		Addr: kafka.TCP("localhost:9092"),   // entry point to our cluster
		Topic: "system-logs",                // target topic we created
		Balancer: &kafka.LeastBytes{},       // algo to decide how message are split between partition
		WriteTimeout: 10 * time.Second,      // max time to wait for the response from the broker
		RequiredAcks: kafka.RequireOne,     // wait for the leader broker to acknowledge receipt
	}

	defer func (){
		if err:= writer.Close(); err !=nil {
			log.Fatalf("failed to close kafka writer: %v", err)
		}
	}()

	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()

	fmt.Println("Starting log generator at exactly 100 RPS")

	for range ticker.C {
		event := models.LogEvent{
			ID: uuid.New().String(),
			TimeStamp: time.Now(),
			Level: levels[rand.Intn(len(levels))],
			Service: services[rand.Intn(len(services))],
			Message: messages[rand.Intn(len(messages))],
		}
		marshalledEvent, marshallErr := json.Marshal(event)
		if marshallErr != nil {
			log.Fatalf("Error occurred during marshalling the event %v", marshallErr)
		} 
		if err := writer.WriteMessages(context.Background(), kafka.Message{Value: marshalledEvent}); err !=nil{
			log.Fatalf("kafka write error %v", err)
		}

		messageCount++

		if messageCount % 100 == 0 {
			fmt.Printf("🚀 Successfully batched and shipped %d logs to Kafka! (Last event: [%s] %s)\n", 
				messageCount, event.Level, event.Service)
		}
	}
}