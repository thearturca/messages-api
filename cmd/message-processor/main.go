package main

import (
	"log"
	messagesProcessor "message-service/internal/messages-processor"
	"os"

	"github.com/kelseyhightower/envconfig"
	"github.com/segmentio/kafka-go"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Kafka struct {
		Brokers []string `envconfig:"KAFKA_BROKERS" yaml:"brokers"`
		Topic   string   `envconfig:"KAFKA_TOPIC" yaml:"topic"`
	} `yaml:"kafka"`
}

var config Config

func init() {
	log.Println("parsing configs...")

	log.Println("loading environment variables...")

	err := envconfig.Process("", &config)

	if err != nil {
		log.Fatalf("failed to load environment variables: %v", err)
	}

	log.Println("loading config file...")

	file, err := os.OpenFile("config.yml", os.O_RDONLY, 0)

	if err != nil {
		log.Printf("failed to open config file: %v\n", err)
		return
	}

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
}

func main() {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: config.Kafka.Brokers,
		Topic:   config.Kafka.Topic,
		GroupID: "message-processor-service",
	})

	writer := &kafka.Writer{
		Addr:                   kafka.TCP(config.Kafka.Brokers...),
		Topic:                  config.Kafka.Topic,
		Balancer:               &kafka.LeastBytes{},
		AllowAutoTopicCreation: true,
	}

	app := messagesProcessor.NewApp(&messagesProcessor.Config{
		Kafka: struct {
			Reader *kafka.Reader
			Writer *kafka.Writer
		}{
			Reader: reader,
			Writer: writer,
		},
	})

	if err := app.Run(); err != nil {
		log.Fatalf("failed to run app: %v", err)
	}
}
