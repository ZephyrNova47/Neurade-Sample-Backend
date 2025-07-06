package entity

type User struct {
	ID           int    `gorm:"column:id;primaryKey"`
	Email        string `gorm:"column:email"`
	PasswordHash string `gorm:"column:password_hash"`
	Role         string `gorm:"column:role"`
	Verified     bool   `gorm:"column:verified"`
	GithubToken  string `gorm:"column:github_token"`
	CreatedAt    string `gorm:"column:created_at"`
	UpdatedAt    string `gorm:"column:updated_at"`
}
