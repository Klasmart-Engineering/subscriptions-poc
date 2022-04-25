package messaging

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"log"
	"subscriptions.demo/database"
	"subscriptions.demo/models"
	"time"
)

//TODO should be completely configurable this entire file should be redone
const (
	TOPIC               = "subscription-evaluation"
	SUBSCRIPTION_CHANGE = "subscription-change"
	PARTITION           = 0
	NETWORK             = "tcp"
	ADDRESS             = "one-node-external.redpanda:9092"
)

func Publish(message string) {
	conn, err := kafka.DialLeader(context.Background(), NETWORK, ADDRESS, TOPIC, PARTITION)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}

	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	_, err = conn.WriteMessages(
		kafka.Message{Value: []byte(message)},
	)
	if err != nil {
		log.Fatal("failed to write messages:", err)
	}

	if err := conn.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}
}

var dbInstance db.Database

func Consume(db db.Database) {
	dbInstance = db
	// make a new reader that consumes from topic-A, partition 0, at offset 42
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{ADDRESS},
		Topic:     SUBSCRIPTION_CHANGE,
		Partition: 0,
		MinBytes:  10e3, // 10KB
		MaxBytes:  10e6, // 10MB
	})

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			break
		}

		var test = string(m.Value)

		data := models.SubscriptionChange{}
		json.Unmarshal([]byte(test), &data)

		dbInstance.UpdateSubscriptionStatus(data.SubscriptionId, data.Active)

	}

	if err := r.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	}

	//conn, err := kafka.DialLeader(context.Background(), NETWORK, ADDRESS, SUBSCRIPTION_CHANGE, PARTITION)
	//if err != nil {
	//	log.Fatal("failed to dial leader:", err)
	//}
	//
	//dbInstance = db
	//conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	//batch := conn.ReadBatch(10e3, 1e6) // fetch 10KB min, 1MB max
	//
	//b := make([]byte, 10e3) // 10KB max per message
	//for {
	//	n, err := batch.Read(b)
	//	if err != nil {
	//		break
	//	}
	//	fmt.Println("CONSUMING")
	//	var test = string(b[:n])
	//	fmt.Println(test)
	//
	//	data := models.SubscriptionChange{}
	//	json.Unmarshal([]byte(test), &data)
	//	dbInstance.UpdateSubscriptionStatus(data.SubscriptionId, data.Active)
	//}
	//
	//if err := batch.Close(); err != nil {
	//	log.Println("failed to close batch:", err)
	//}
	//
	//if err := conn.Close(); err != nil {
	//	log.Println("failed to close connection:", err)
	//}
}
