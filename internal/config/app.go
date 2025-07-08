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
	courseRepo := repository.NewCourseRepository(config.DB, config.Log)
	asisgnmentRepo := repository.NewAssignmentRepository(config.DB, config.Log)
	prRepo := repository.NewPrRepository(config.DB, config.Log)

	userService := service.NewUserService(config.DB, userRepo, config.Log)
	llmService := service.NewLLMService(config.DB, llmRepo, config.Log)
	courseService := service.NewCourseService(config.DB, courseRepo, config.Log)
	assignmentService := service.NewAssignmentService(config.DB, asisgnmentRepo, config.Log)
	prService := service.NewPrService(config.DB, prRepo, config.Log)
	githubService := service.NewGitHubService(config.Log)

	userController := controller.NewUserController(userService, config.Log, config.Config.JWTSecret)
	llmController := controller.NewLLMController(llmService, config.Log)
	courseController := controller.NewCourseController(courseService, config.Log, config.Minio)
	assignmentController := controller.NewAssignmentController(assignmentService, config.Log, config.Minio)
	prController := controller.NewPrController(prService, config.Log)
	githubWebhookController := controller.NewGitHubWebhookController(githubService, prService, courseService, config.Log)

	r := route.RouteConfig{
		App:                     chi.NewRouter(),
		UserController:          userController,
		LLMController:           llmController,
		CourseController:        courseController,
		AssignmentController:    assignmentController,
		PrController:            prController,
		GitHubWebhookController: githubWebhookController,
	}

	return r.Setup()
}
