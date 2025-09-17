package consumer

import (
	"context"
	"encoding/json"
	"log"

	"github.com/iskhakmuhamad/mylogs-ws/internal/entity"
	"github.com/iskhakmuhamad/mylogs-ws/internal/usecase"
	"github.com/iskhakmuhamad/mylogs-ws/pkg/broker"
	"github.com/segmentio/kafka-go"
)

type NotificationConsumer struct {
	consumer *broker.KafkaConsumer
	usecase  usecase.NotificationUsecase
}

func NewNotificationConsumer(c *broker.KafkaConsumer, u usecase.NotificationUsecase) *NotificationConsumer {
	return &NotificationConsumer{consumer: c, usecase: u}
}

func (nc *NotificationConsumer) Start(ctx context.Context) error {
	return nc.consumer.Consume(ctx, func(m kafka.Message) error {
		var notif entity.Notification
		if err := json.Unmarshal(m.Value, &notif); err != nil {
			log.Printf("❌ Failed to unmarshal message: %v", err)
			return err
		}

		// call business logic
		if err := nc.usecase.HandleIncomingNotification(ctx, notif); err != nil {
			log.Printf("❌ Usecase failed: %v", err)
			return err
		}
		log.Printf("✅ Processed notification: %+v", notif)
		return nil
	})
}
