package controller

import (
	"be/neurade/v2/internal/model/converter"
	"be/neurade/v2/internal/service"
	"be/neurade/v2/internal/util"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/minio/minio-go/v7"
	"github.com/sirupsen/logrus"
)

type CourseController struct {
	CourseService *service.CourseService
	Log           *logrus.Logger
	MinioUntil    *util.MinioUtil
}

func NewCourseController(courseService *service.CourseService, log *logrus.Logger, minioClient *minio.Client) *CourseController {
	return &CourseController{CourseService: courseService, Log: log,
		MinioUntil: util.NewMinioUtil(minioClient, log)}
}

func (c *CourseController) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(64 << 20); err != nil {
		c.Log.WithError(err).Error("Error parsing multipart form")
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	course := converter.RequestToCourseRequest(r)
	var generalAnswerContent string
	generalAnswerFile, fileHeader, err := r.FormFile("general_answer")
	if err == nil {
		defer generalAnswerFile.Close()
		c.Log.Infof("Uploaded file: %s, size: %d bytes", fileHeader.Filename, fileHeader.Size)

		content, err := io.ReadAll(generalAnswerFile)
		if err != nil {
			c.Log.WithError(err).Error("Error reading general answer file")
			http.Error(w, "Error reading file", http.StatusInternalServerError)
			return
		}
		generalAnswerContent = string(content)
	}
	generalAnswerURL, _ := c.MinioUntil.SaveFile(r.Context(), course.CourseName, course.CreatedAt, "course", fileHeader.Filename, generalAnswerContent)
	course.GeneralAnswer = generalAnswerURL
	courseResponse, err := c.CourseService.Create(r.Context(), course)
	if err != nil {
		c.Log.Println("Failed to create course:", err)
		http.Error(w, "Failed to create course", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(courseResponse)
}

func (c *CourseController) GetByID(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "course_id"))
	courseResponse, err := c.CourseService.GetByID(r.Context(), id)
	generalAnsPres, _ := c.MinioUntil.GetFile(r.Context(), courseResponse.GeneralAnswer)
	courseResponse.GeneralAnswer = generalAnsPres
	if err != nil {
		c.Log.Println("Failed to get course:", err)
		http.Error(w, "Failed to get course", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(courseResponse)
}

func (c *CourseController) GetAllByOwner(w http.ResponseWriter, r *http.Request) {
	ownerID, _ := strconv.Atoi(chi.URLParam(r, "user_id"))
	courseResponse, err := c.CourseService.GetAllByOwner(r.Context(), ownerID)
	if err != nil {
		c.Log.Println("Failed to get all course:", err)
		http.Error(w, "Failed to get all course", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(courseResponse)
}
