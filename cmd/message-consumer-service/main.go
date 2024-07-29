package main

import (
	"log"
	consumer "message-service/internal/consumer"
	"message-service/internal/db"
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

	DB struct {
		PG_HOST     string `envconfig:"PG_HOST" yaml:"host"`
		PG_PORT     string `envconfig:"PG_PORT" yaml:"port"`
		PG_USER     string `envconfig:"PG_USER" yaml:"user"`
		PG_PASSWORD string `envconfig:"PG_PASSWORD" yaml:"password"`
		PG_DATABASE string `envconfig:"PG_DATABASE" yaml:"database"`
	} `yaml:"db"`
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
		log.Fatalf("failed to open config file: %v", err)
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
	log.Println("connecting to db...")

	db, err := db.New(&db.Config{
		Host:     config.DB.PG_HOST,
		Port:     config.DB.PG_PORT,
		User:     config.DB.PG_USER,
		Password: config.DB.PG_PASSWORD,
		Database: config.DB.PG_DATABASE,
	})

	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}

	log.Println("db connected")

	kafka := kafka.NewReader(kafka.ReaderConfig{
		Brokers: config.Kafka.Brokers,
		Topic:   config.Kafka.Topic,
		GroupID: "message-consumer-service",
	})

	app := consumer.NewApp(&consumer.Config{
		DB:    db,
		Kafka: kafka,
	})

	if err := app.Run(); err != nil {
		log.Fatalf("failed to run app: %v", err)
	}
}
