package repository

import (
	"be/neurade/v2/internal/entity"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type CourseRepository struct {
	Repository[entity.Course]
	Log *logrus.Logger
}

func NewCourseRepository(db *gorm.DB, log *logrus.Logger) *CourseRepository {
	return &CourseRepository{
		Repository: Repository[entity.Course]{
			DB: db,
		},
		Log: log,
	}
}

func (r *CourseRepository) FindAllByOwner(db *gorm.DB, courses *[]entity.Course, ownerID int) error {
	return db.Where("user_id = ?", ownerID).Find(courses).Error
}

func (r *CourseRepository) FindAllByGithubURL(db *gorm.DB, course *entity.Course, githubURL string) error {
	return db.Where("github_url = ?", githubURL).Find(course).Error
}
