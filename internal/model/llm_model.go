package model

type LLMCreateRequest struct {
	UserID     int    `json:"user_id"`
	ModelName  string `json:"model_name" validate:"required"`
	ModelToken string `json:"model_token" validate:"required"`
	Status     string `json:"status" validate:"required,oneof=active inactive"`
}

type LLMResponse struct {
	ID         int    `json:"id"`
	UserID     int    `json:"user_id"`
	ModelName  string `json:"model_name"`
	ModelToken string `json:"model_token"`
	Status     string `json:"status"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}
