package controller

import (
	"be/neurade/v2/internal/model"
	"be/neurade/v2/internal/service"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

type GitHubWebhookController struct {
	githubService *service.GitHubService
	prService     *service.PrService
	courseService *service.CourseService
	log           *logrus.Logger
}

func NewGitHubWebhookController(githubService *service.GitHubService, prService *service.PrService, courseService *service.CourseService, log *logrus.Logger) *GitHubWebhookController {
	return &GitHubWebhookController{
		githubService: githubService,
		prService:     prService,
		courseService: courseService,
		log:           log,
	}
}

func (c *GitHubWebhookController) FetchPullRequests(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		c.log.Println("Failed to parse multipart form:", err)
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}
	courseID, _ := strconv.Atoi(r.FormValue("course_id"))
	req := &model.FetchPullRequestsRequest{
		CourseID:    courseID,
		GithubURL:   r.FormValue("github_url"),
		GithubToken: r.FormValue("github_token"),
	}
	c.log.Info(r.FormValue("github_url"))

	// _, err := c.courseService.GetByID(r.Context(), req.CourseID)
	// if err != nil {
	// 	c.log.Errorf("Course not found: %v", err)
	// 	http.Error(w, "Course not found", http.StatusNotFound)
	// 	return
	// }

	pullRequests, err := c.githubService.GetPullRequests(r.Context(), req.GithubURL, req.GithubToken)
	if err != nil {
		c.log.Errorf("Failed to fetch pull requests: %v", err)
		http.Error(w, "Failed to fetch pull requests from GitHub", http.StatusInternalServerError)
		return
	}

	// Save pull requests to database
	savedCount := 0
	for _, pr := range pullRequests {
		prRequest := &model.PrCreateRequest{
			CourseID:      req.CourseID,
			PrName:        pr.Title,
			PrDescription: pr.Body,
			PrNumber:      pr.Number,
			Status:        pr.State,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}

		_, err := c.prService.Create(r.Context(), prRequest)
		if err != nil {
			c.log.Errorf("Failed to save PR %d: %v", pr.Number, err)
			continue
		}
		savedCount++
	}

	// Return response
	response := model.FetchPullRequestsResponse{
		Message:           "Successfully fetched and saved pull requests",
		PullRequestsCount: savedCount,
		CourseID:          req.CourseID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

	c.log.Infof("Successfully processed %d pull requests for course %d", savedCount, req.CourseID)
}

func (c *GitHubWebhookController) GetPullRequestsByCourse(w http.ResponseWriter, r *http.Request) {
	courseID, err := strconv.Atoi(chi.URLParam(r, "course_id"))
	if err != nil {
		http.Error(w, "Invalid course ID", http.StatusBadRequest)
		return
	}

	course, err := c.courseService.GetByID(r.Context(), courseID)
	if err != nil {
		http.Error(w, "Course not found", http.StatusNotFound)
		return
	}
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		c.log.Println("Failed to parse multipart form:", err)
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}
	githubToken := r.FormValue("github_token")

	if githubToken == "" {
		http.Error(w, "GitHub token not found", http.StatusBadRequest)
		return
	}

	pullRequests, err := c.githubService.GetPullRequests(r.Context(), course.GithubURL, githubToken)
	if err != nil {
		c.log.Errorf("Failed to fetch pull requests: %v", err)
		http.Error(w, "Failed to fetch pull requests", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pullRequests)
}
