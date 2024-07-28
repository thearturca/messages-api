package messagesProcessor

import (
	"context"
	"encoding/json"
	"log"
	"message-service/internal/db"
	"time"

	"github.com/segmentio/kafka-go"
)

type Service struct {
	KafkaWriter *kafka.Writer
}

func NewService(writer *kafka.Writer) *Service {
	return &Service{
		KafkaWriter: writer,
	}
}

func (s *Service) Handle(ctx context.Context, message kafka.Message) error {
	var dto db.Message
	json.Unmarshal(message.Value, &dto)

	if dto.IsProcessed {
		return nil
	}

	processedMessage := s.ProcessMessage(context.Background(), dto)

	if marshalledMessage, err := json.Marshal(processedMessage); err != nil {
		log.Println("failed to marshal message: ", err)
		return err
	} else {
		s.KafkaWriter.WriteMessages(context.Background(), kafka.Message{
			Key:   []byte(processedMessage.Id),
			Value: marshalledMessage,
		})
	}

	return nil
}

func (s *Service) ProcessMessage(ctx context.Context, message db.Message) *db.Message {
	log.Println("process message: ", message.Id)
	log.Println("message text: ", message.Text)

	time.Sleep(20 * time.Second)

	log.Println("message processed: ", message.Id)

	now := time.Now()

	message.ProcessedAt = &now
	message.IsProcessed = true

	return &message
}
