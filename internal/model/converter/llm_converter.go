package converter

import (
	"be/neurade/v2/internal/entity"
	"be/neurade/v2/internal/model"
)

func LLMToResponse(llm *entity.LLM) *model.LLMResponse {
	return &model.LLMResponse{
		ID:         llm.ID,
		UserID:     llm.UserID,
		ModelName:  llm.ModelName,
		ModelToken: llm.ModelToken,
		Status:     llm.Status,
		CreatedAt:  llm.CreatedAt,
		UpdatedAt:  llm.UpdatedAt,
	}
}
