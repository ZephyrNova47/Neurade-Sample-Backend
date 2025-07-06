package service

import (
	"be/neurade/v2/internal/entity"
	"be/neurade/v2/internal/model"
	"be/neurade/v2/internal/model/converter"
	"be/neurade/v2/internal/repository"
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type LLMService struct {
	DB            *gorm.DB
	LLMRepository *repository.LLMRepository
	Log           *logrus.Logger
}

func NewLLMService(db *gorm.DB, llmRepository *repository.LLMRepository, log *logrus.Logger) *LLMService {
	return &LLMService{
		DB:            db,
		LLMRepository: llmRepository,
		Log:           log,
	}
}

func (s *LLMService) Create(ctx context.Context, request *model.LLMCreateRequest) (*model.LLMResponse, error) {

	tx := s.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			s.Log.Error("Transaction failed: ", r)
		}
	}()

	llm := &entity.LLM{
		UserID:     request.UserID,
		ModelName:  request.ModelName,
		ModelToken: request.ModelToken,
		Status:     request.Status,
		CreatedAt:  time.Now().Format(time.RFC3339),
		UpdatedAt:  time.Now().Format(time.RFC3339),
	}

	if err := tx.Create(llm).Error; err != nil {
		tx.Rollback()
		s.Log.Error("Failed to create LLM: ", err)
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		s.Log.Error("Failed to commit transaction: ", err)
		return nil, err
	}

	return converter.LLMToResponse(llm), nil
}

func (s *LLMService) GetByID(ctx context.Context, id int) (*model.LLMResponse, error) {
	tx := s.DB.WithContext(ctx).Begin()
	llm := entity.LLM{}

	err := s.LLMRepository.FindById(tx, &llm, id)
	if err != nil {
		s.Log.Error("Failed to find LLM by ID: ", err)
		return nil, err
	}
	return converter.LLMToResponse(&llm), nil
}
