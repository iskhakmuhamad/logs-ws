package repository

import (
	"github.com/iskhakmuhamad/mylogs-ws/internal/entity"
	"gorm.io/gorm"
)

type NotificationRepository interface {
	Create(n *entity.Notification) error
	FindByUser(userID uint, unreadOnly bool) ([]entity.Notification, error)
	MarkAsRead(id uint) error
}

type notificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) NotificationRepository {
	return &notificationRepository{db: db}
}

func (r *notificationRepository) Create(n *entity.Notification) error {
	return r.db.Create(n).Error
}

func (r *notificationRepository) FindByUser(userID uint, unreadOnly bool) ([]entity.Notification, error) {
	var list []entity.Notification
	q := r.db.Where("user_id = ?", userID).Order("created_at DESC")
	if unreadOnly {
		q = q.Where("read = false")
	}
	err := q.Find(&list).Error
	return list, err
}

func (r *notificationRepository) MarkAsRead(id uint) error {
	return r.db.Model(&entity.Notification{}).Where("id = ?", id).Update("read", true).Error
}
