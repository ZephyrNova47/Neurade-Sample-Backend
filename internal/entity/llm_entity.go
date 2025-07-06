package entity

type LLM struct {
	ID         int    `gorm:"column:id;primaryKey"`
	UserID     int    `gorm:"column:user_id"`
	ModelName  string `gorm:"column:model_name"`
	ModelToken string `gorm:"column:model_token"`
	Status     string `gorm:"column:status"`
	CreatedAt  string `gorm:"column:created_at"`
	UpdatedAt  string `gorm:"column:updated_at"`
}
