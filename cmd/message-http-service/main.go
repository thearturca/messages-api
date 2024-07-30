package main

import (
	"log"
	"message-service/internal/db"
	"message-service/internal/messages"
	"os"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/segmentio/kafka-go"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Kafka struct {
		Brokers []string `envconfig:"KAFKA_BROKERS" yaml:"brokers"`
		Topic   string   `envconfig:"KAFKA_TOPIC" yaml:"topic"`
	} `yaml:"kafka"`

	Server struct {
		PORT string `yaml:"port" envconfig:"PORT"`
		HOST string `yaml:"host" envconfig:"HOST"`
		Auth struct {
			Username string `yaml:"username" envconfig:"BASIC_AUTH_USERNAME"`
			Password string `yaml:"password" envconfig:"BASIC_AUTH_PASSWORD"`
		} `yaml:"auth"`
	} `yaml:"server"`

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
	log.Println("connecting to db...")

	db, err := db.New(&db.Config{
		Host:     config.DB.PG_HOST,
		Port:     config.DB.PG_PORT,
		User:     config.DB.PG_USER,
		Password: config.DB.PG_PASSWORD,
		Database: config.DB.PG_DATABASE,
	})

	defer db.Close()

	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}

	log.Println("db connected")

	kafkaWriter := &kafka.Writer{
		Topic:                  config.Kafka.Topic,
		Addr:                   kafka.TCP(config.Kafka.Brokers...),
		Balancer:               &kafka.LeastBytes{},
		AllowAutoTopicCreation: true,
		RequiredAcks:           kafka.RequireNone,
		WriteBackoffMin:        10 * time.Millisecond,
	}

	defer kafkaWriter.Close()

	app := messages.NewApp(&messages.Config{
		Port: config.Server.PORT,
		Host: config.Server.HOST,
		Auth: messages.Auth{
			Username: config.Server.Auth.Username,
			Password: config.Server.Auth.Password,
		},
		DB:    db,
		Kafka: kafkaWriter,
	})

	log.Printf("starting server on port %s", config.Server.PORT)

	err = app.Run()

	if err != nil {
		log.Fatalf("failed to run app: %s", err)
	}
}
