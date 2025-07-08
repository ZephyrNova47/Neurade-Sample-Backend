package service

import (
	"be/neurade/v2/internal/entity"
	"be/neurade/v2/internal/model"
	"be/neurade/v2/internal/model/converter"
	"be/neurade/v2/internal/repository"
	"context"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type CourseService struct {
	DB               *gorm.DB
	CourseRepository *repository.CourseRepository
	Log              *logrus.Logger
}

func NewCourseService(db *gorm.DB, courseRepository *repository.CourseRepository, log *logrus.Logger) *CourseService {
	return &CourseService{DB: db, CourseRepository: courseRepository, Log: log}
}

func (s *CourseService) Create(ctx context.Context, request *model.CourseCreateRequest) (*model.CourseResponse, error) {
	tx := s.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	course := converter.CourseToEntity(request)

	if err := tx.Create(course).Error; err != nil {
		s.Log.WithContext(ctx).WithError(err).Error("failed to create course")
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		s.Log.WithContext(ctx).WithError(err).Error("failed to commit transaction")
		return nil, err
	}

	return converter.CourseToResponse(course), nil
}

func (s *CourseService) GetByID(ctx context.Context, id int) (*model.CourseResponse, error) {
	course := &entity.Course{}
	err := s.CourseRepository.FindById(s.DB, course, id)
	if err != nil {
		s.Log.WithContext(ctx).WithError(err).Error("failed to get course by id")
		return nil, err
	}
	return converter.CourseToResponse(course), nil
}

func (s *CourseService) GetByGithubURL(ctx context.Context, githubURL string) (*model.CourseResponse, error) {
	course := &entity.Course{}
	err := s.CourseRepository.FindAllByGithubURL(s.DB, course, githubURL)
	if err != nil {
		s.Log.WithContext(ctx).WithError(err).Error("failed to get course by id")
		return nil, err
	}
	return converter.CourseToResponse(course), nil
}

func (s *CourseService) GetAllByOwner(ctx context.Context, ownerID int) ([]*model.CourseResponse, error) {
	courses := make([]entity.Course, 0)
	err := s.CourseRepository.FindAllByOwner(s.DB, &courses, ownerID)
	if err != nil {
		s.Log.WithContext(ctx).WithError(err).Error("failed to get all course")
		return nil, err
	}

	courseResponses := make([]*model.CourseResponse, 0)
	for i := range courses {
		courseResponses = append(courseResponses, converter.CourseToResponse(&courses[i]))
	}
	return courseResponses, nil
}
