package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"realtime-log-aggregation/internal/db"
	"realtime-log-aggregation/internal/models"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/segmentio/kafka-go"
)

const batchSizeLimit = 50
const maxWaitTime = 500 * time.Millisecond

func main(){
    // load the env file on startup
	// this pull the variables out of the text file and inject them into OS environment
	if envLoadErr := godotenv.Load(); envLoadErr !=nil{
		log.Println("No .env file found, relying on system environment variables")
	}
	reader:= kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic: "system-logs",
		GroupID: "log-processor-group",
		MinBytes: 1,
		MaxBytes: 10e6,
	})

	defer func(){
		if err := reader.Close(); err != nil {
			log.Fatalf("Error during closing kafka reader %v", err)
		}
	}()

	connString := os.Getenv("DATABASE_URL")
	if connString == ""{
		log.Fatalf("Critical configuration error: DATABASE_URL environment variable is not set")
	}

	pgxPool, dbErr := db.InitDB(connString)
	if dbErr != nil {
		log.Fatalf("Failed to initialize database %v", dbErr)
	}
	defer pgxPool.Close()

	logBatch := make([]models.LogEvent, 0, batchSizeLimit)
	kafkaBatch := make([]kafka.Message, 0, batchSizeLimit)
	lastFlush := time.Now()

	for {
		// using timeout context so that loop perform a flush if no logs have arrived in a while
		ctx, cancel := context.WithTimeout(context.Background(), 100 * time.Millisecond)
		// Read message from kafka topic 
		//We are using fetch message instead of read message because we want to commit the offset after processing the message
		kafkaMsg, readErr :=reader.FetchMessage(ctx)
		cancel()
		if readErr != nil && readErr != context.DeadlineExceeded {
			log.Printf("Error during reading kafka message %v", readErr)
			continue
		}
		timeSinceLastFlush := time.Since(lastFlush)
		shouldFlush := len(logBatch) >= batchSizeLimit || (len(logBatch) > 0 && timeSinceLastFlush >= maxWaitTime)
		if shouldFlush {
			rowsInserted, dbErr := dbPush(pgxPool, logBatch)
			if dbErr != nil {
				log.Fatalf("Db insertion failure %v", dbErr)
			} 
			log.Printf("Successfully flushed logs to postgres, row inserted: %v", rowsInserted)
			if commitErr := reader.CommitMessages(context.Background(), kafkaBatch...); commitErr != nil{
				log.Printf("Error during committing the batch offset %v", commitErr)
			}
			lastFlush = time.Now()
			// Helps in optimizing garbage collection pause times as reset length to zero so go overwrites existing slots instead of asking cpu to allocate new memory
			logBatch =logBatch[:0]
			kafkaBatch = kafkaBatch[:0]
		}
		//Handle Timeout Check: Skip unmarshalling if no message was fetched
		if readErr == context.DeadlineExceeded {
				continue
		}
		var event = models.LogEvent{}
		if unmarshallError := json.Unmarshal(kafkaMsg.Value, &event); unmarshallError != nil{
			log.Printf("Error during unmarshalling message %v", unmarshallError)
			continue
		}
		logBatch = append(logBatch, event)
		kafkaBatch = append(kafkaBatch, kafkaMsg)
		fmt.Println("Message read successfully ")
	}
}

func dbPush(pgxPool *pgxpool.Pool ,logBatch []models.LogEvent) (int64, error) {
	rows := make([][]any, len(logBatch))
	columns := []string{"id", "timestamp", "level", "service", "message"}
	for i, logItem := range logBatch {
		rows[i] = []any{logItem.ID, logItem.TimeStamp, logItem.Level, logItem.Service, logItem.Message}
	}
	rowInserted, dbInsertErr := pgxPool.CopyFrom(context.Background(), pgx.Identifier{"system_logs"}, columns, pgx.CopyFromRows(rows))
	if dbInsertErr != nil {
		return rowInserted, fmt.Errorf("Error during row insertion %v", dbInsertErr)
	}
	return rowInserted, nil
}