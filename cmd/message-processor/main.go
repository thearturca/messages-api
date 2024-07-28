package main

import (
	"context"
	"encoding/json"
	"log"
	"message-service/internal/db"
	"os"
	"time"

	"github.com/segmentio/kafka-go"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Kafka struct {
		Brokers []string `envconfig:"KAFKA_BROKERS" yaml:"brokers"`
		Topic   string   `envconfig:"KAFKA_TOPIC" yaml:"topic"`
	} `yaml:"kafka"`
}

func main() {
	log.Println("parsing config...")
	file, err := os.OpenFile("config.yml", os.O_RDONLY, 0)

	if err != nil {
		log.Fatalf("failed to open config file: %v", err)
	}

	var config Config
	decoder := yaml.NewDecoder(file)

	err = decoder.Decode(&config)

	if err != nil {
		log.Fatalf("failed to decode config file: %v", err)
	}

	err = file.Close()

	if err != nil {
		log.Fatalf("failed to close config file: %v", err)
	}

	log.Println("config parsed")

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: config.Kafka.Brokers,
		Topic:   config.Kafka.Topic,
	})

	writer := kafka.Writer{
		Addr:                   kafka.TCP(config.Kafka.Brokers...),
		Topic:                  config.Kafka.Topic,
		Balancer:               &kafka.LeastBytes{},
		AllowAutoTopicCreation: true,
	}

	for {
		msg, err := reader.ReadMessage(context.Background())

		if err != nil {
			break
		}

		log.Println(msg)

		var dtoMessage *db.Message

		err = json.Unmarshal(msg.Value, &dtoMessage)

		if err != nil {
			log.Fatal(err)
		}

		log.Println(dtoMessage)
		time.Sleep(5 * time.Second)
		writer.WriteMessages(context.Background(), kafka.Message{Key: []byte(dtoMessage.Id), Value: dtoMessage})

	}

	if err := reader.Close(); err != nil {
		log.Fatal(err)
	}

}
