package entity

import "time"

type Notification struct {
	ID        uint      `json:"id" gorm:"primaryKey;column:id;autoIncrement"`
	UserID    uint      `json:"user_id" gorm:"column:user_id;index;not null"`
	Title     string    `json:"title" gorm:"column:title;type:varchar(255);not null"`
	Body      string    `json:"body" gorm:"column:body;type:text"`
	Read      bool      `json:"read" gorm:"column:read;default:false;not null"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
}

func (Notification) TableName() string {
	return "notifications"
}
