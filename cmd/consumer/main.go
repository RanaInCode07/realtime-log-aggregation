package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"realtime-log-aggregation/internal/models"

	"github.com/segmentio/kafka-go"
)

func main(){
	reader:= kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic: "system-logs",
		GroupID: "log-processor-group",
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})

	defer func(){
		if err := reader.Close(); err != nil {
			log.Fatalf("Error during closing kafka reader %v", err)
		}
	}()

	for {
		msg, readErr :=reader.ReadMessage(context.Background())      // it is a blocking call wait until a log lands in kafka
		if readErr != nil {
			log.Printf("Error during reading kafka message %v", readErr)
			continue
		}
		var event = models.LogEvent{}
		if unmarshallError := json.Unmarshal(msg.Value, &event); unmarshallError != nil{
			log.Printf("Error during unmarshalling message %v", unmarshallError)
			continue
		}
		fmt.Println("Message read successfully ")
	}
}