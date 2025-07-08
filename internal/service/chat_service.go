package service

import (
	"be/neurade/v2/internal/repository"

	"github.com/sirupsen/logrus"
)

type ChatService struct {
	ChatRepository *repository.ChatRepository
	Log            *logrus.Logger
}

func NewChatRepository(chatRepository *repository.ChatRepository, log *logrus.Logger) *ChatService {
	return &ChatService{
		ChatRepository: chatRepository,
		Log:            log,
	}
}
