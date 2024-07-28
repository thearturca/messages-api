package consumer

import (
	"context"
	"encoding/json"
	"log"
	"message-service/internal/db"
	"message-service/internal/messages"
	"time"

	"github.com/segmentio/kafka-go"
)

type Service struct {
	messagesService *messages.Service
}

func NewService(messagesService *messages.Service) *Service {
	return &Service{
		messagesService: messagesService,
	}
}

func (s *Service) Handle(ctx context.Context, message kafka.Message) error {
	var dto db.Message
	json.Unmarshal(message.Value, &dto)

	if !dto.IsProcessed {
		return nil
	}

	_, err := s.ProcessMessage(context.Background(), dto)

	if err != nil {
		log.Println("failed to process message: ", err)
		return err
	}

	return nil
}

func (s *Service) ProcessMessage(ctx context.Context, message db.Message) (*db.Message, error) {
	if message.ProcessedAt == nil {
		now := time.Now()
		message.ProcessedAt = &now
	}

	updatedMessage, err := s.messagesService.UpdateMessage(ctx, &db.Message{
		Id:          message.Id,
		IsProcessed: true,
		ProcessedAt: message.ProcessedAt,
	})

	if err != nil {
		return nil, err
	}

	return updatedMessage, nil
}
