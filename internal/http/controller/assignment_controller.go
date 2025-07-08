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

type AssignmentController struct {
	AssignmentService *service.AssignmentService
	Log               *logrus.Logger
	MinioUtil         *util.MinioUtil
}

func NewAssignmentController(assignmentService *service.AssignmentService, log *logrus.Logger, minioClient *minio.Client) *AssignmentController {
	return &AssignmentController{AssignmentService: assignmentService, Log: log,
		MinioUtil: util.NewMinioUtil(minioClient, log)}
}

func (c *AssignmentController) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(64 << 20); err != nil {
		c.Log.Println("Failed to parse multipart form:", err)
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	assignment := converter.RequestToAssignmentRequest(r)
	var assignmentContent string
	assignmentFile, fileHeader, err := r.FormFile("assignment_file")
	if err == nil {
		defer assignmentFile.Close()
		c.Log.Infof("Upload file: %s, size: %d bytes", fileHeader.Filename, fileHeader.Size)

		content, err := io.ReadAll(assignmentFile)
		if err != nil {
			c.Log.WithError(err).Error("Error reading assignment file")
			http.Error(w, "Error reading assignment file", http.StatusInternalServerError)
			return
		}
		assignmentContent = string(content)
	}
	course, _ := c.AssignmentService.GetCourseByID(r.Context(), assignment.CourseID)
	assignmentURL, _ := c.MinioUtil.SaveFile(r.Context(), course.CourseName, course.CreatedAt, "assignment", assignment.AssignmentName, assignmentContent)
	assignment.AssignmentURL = assignmentURL
	assignmentResponse, err := c.AssignmentService.Create(r.Context(), assignment)
	if err != nil {
		c.Log.Println("Failed to create assignment:", err)
		http.Error(w, "Failed to create assignment:", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(assignmentResponse)
}

func (c *AssignmentController) GetByID(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "assignment_id"))
	assignmentResponse, err := c.AssignmentService.GetByID(r.Context(), id)
	assignmentPres, _ := c.MinioUtil.GetFile(r.Context(), assignmentResponse.AssignmentURL)
	assignmentResponse.AssignmentURL = assignmentPres
	if err != nil {
		c.Log.Println("Failed to get assignment by ID")
		http.Error(w, "Failed to get assignment by ID", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(assignmentResponse)
}

func (c *AssignmentController) GetAllByCourse(w http.ResponseWriter, r *http.Request) {
	courseID, _ := strconv.Atoi(chi.URLParam(r, "course_id"))
	assignmentResponse, err := c.AssignmentService.GetAllByCourse(r.Context(), courseID)
	if err != nil {
		c.Log.Println("Failed to get all assignment by course")
		http.Error(w, "Failed to get all assignment by course", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(assignmentResponse)
}
