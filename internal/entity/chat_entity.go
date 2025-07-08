package entity

import (
	"encoding/json"
	"time"
)

type Chat struct {
	ID          int             `gorm:"column:id;primaryKey"`
	PrID        int             `gorm:"column:pr_id"`
	ChatHistory json.RawMessage `gorm:"column:chat_history"`
	CreatedAt   time.Time       `gorm:"column:created_at"`
	UpdatedAt   time.Time       `gorm:"column:updated_at"`
}
