package kafkaService

import (
	"context"
	"encoding/json"
	"message-service/internal/db"

	"github.com/segmentio/kafka-go"
)

type WriterService struct {
	kafka *kafka.Writer
}

func NewWriterService(kafka *kafka.Writer) *WriterService {
	return &WriterService{
		kafka: kafka,
	}
}

func (s *WriterService) SendMessage(ctx context.Context, message *db.Message) error {
	marshalledMessage, err := json.Marshal(message)

	if err != nil {
		return err
	}

	return s.kafka.WriteMessages(ctx, kafka.Message{
		Key:   []byte(message.Id),
		Value: marshalledMessage,
	})
}
