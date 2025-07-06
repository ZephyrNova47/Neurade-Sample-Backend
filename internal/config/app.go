package config

import (
	"be/neurade/v2/internal/http/controller"
	"be/neurade/v2/internal/http/route"
	"be/neurade/v2/internal/repository"
	"be/neurade/v2/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/minio/minio-go/v7"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
	DB        *gorm.DB
	Log       *logrus.Logger
	Agent     *AgentConfig
	Minio     *minio.Client
	Config    *Config
	JWTConfig *JWTConfig
}

func Bootstrap(config *BootstrapConfig) *chi.Mux {
	userRepo := repository.NewUserRepository(config.DB, config.Log)
	llmRepo := repository.NewLLMRepository(config.DB, config.Log)
	userService := service.NewUserService(config.DB, userRepo, config.Log)
	llmService := service.NewLLMService(config.DB, llmRepo, config.Log)

	userController := controller.NewUserController(userService, config.Log, config.Config.JWTSecret)
	llmController := controller.NewLLMController(llmService, config.Log)

	r := route.RouteConfig{
		App:            chi.NewRouter(),
		UserController: userController,
		LLMController:  llmController,
	}

	return r.Setup()
}
