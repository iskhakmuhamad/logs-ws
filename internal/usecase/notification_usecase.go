package usecase

import (
	"context"
	"log"

	"github.com/iskhakmuhamad/mylogs-ws/internal/delivery/ws"
	"github.com/iskhakmuhamad/mylogs-ws/internal/entity"
	"github.com/iskhakmuhamad/mylogs-ws/internal/repository"
)

type NotificationUsecase interface {
	Create(n *entity.Notification) error
	List(userID uint, unreadOnly bool) ([]entity.Notification, error)
	MarkAsRead(id uint) error
	HandleIncomingNotification(ctx context.Context, notif entity.Notification) error
}

type notificationUsecase struct {
	repo repository.NotificationRepository
	hub  *ws.Hub
}

func NewNotificationUsecase(r repository.NotificationRepository, h *ws.Hub) NotificationUsecase {
	return &notificationUsecase{repo: r, hub: h}
}

func (u *notificationUsecase) Create(n *entity.Notification) error {
	if err := u.repo.Create(n); err != nil {
		return err
	}

	u.hub.Broadcast(n.UserID, map[string]any{
		"type": "notification",
		"data": n,
	})

	return nil
}

func (u *notificationUsecase) HandleIncomingNotification(ctx context.Context, notif entity.Notification) error {
	log.Println("check ,", notif)
	if err := u.repo.Create(&notif); err != nil {
		return err
	}

	u.hub.Broadcast(notif.UserID, map[string]any{
		"type": "notification",
		"data": notif,
	})

	return nil
}

func (u *notificationUsecase) List(userID uint, unreadOnly bool) ([]entity.Notification, error) {
	return u.repo.FindByUser(userID, unreadOnly)
}

func (u *notificationUsecase) MarkAsRead(id uint) error {
	return u.repo.MarkAsRead(id)
}
