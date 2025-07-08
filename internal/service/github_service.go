package service

import (
	"be/neurade/v2/internal/model"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
)

type GitHubService struct {
	Log *logrus.Logger
}

func NewGitHubService(log *logrus.Logger) *GitHubService {
	return &GitHubService{
		Log: log,
	}
}

func (s *GitHubService) GetPullRequests(ctx context.Context, githubURL string, githubToken string) ([]model.GitHubPullRequest, error) {
	// Extract owner and repo from GitHub URL
	// Example: https://github.com/owner/repo -> owner/repo
	parts := strings.Split(strings.TrimPrefix(githubURL, "https://github.com/"), "/")
	// if len(parts) != 2 {
	// 	return nil, fmt.Errorf("invalid GitHub URL format: %s", githubURL)
	// }
	s.Log.Info(githubURL)
	owner := parts[0]
	repo := parts[1]

	// GitHub API endpoint for pull requests
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/pulls?state=all", owner, repo)

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Authorization", "token "+githubToken)
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", "Neurade-Backend")

	// Make request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("GitHub API error: %s - %s", resp.Status, string(body))
	}

	// Parse response
	var pullRequests []model.GitHubPullRequest
	if err := json.NewDecoder(resp.Body).Decode(&pullRequests); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	s.Log.Infof("Fetched %d pull requests from %s/%s", len(pullRequests), owner, repo)
	return pullRequests, nil
}

func (s *GitHubService) GetRepositoryInfo(ctx context.Context, githubURL, githubToken string) (*model.GitHubRepository, error) {
	// Extract owner and repo from GitHub URL
	parts := strings.Split(strings.TrimPrefix(githubURL, "https://github.com/"), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid GitHub URL format: %s", githubURL)
	}

	owner := parts[0]
	repo := parts[1]

	// GitHub API endpoint for repository info
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s", owner, repo)

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Authorization", "token "+githubToken)
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", "Neurade-Backend")

	// Make request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("GitHub API error: %s - %s", resp.Status, string(body))
	}

	// Parse response
	var repository model.GitHubRepository
	if err := json.NewDecoder(resp.Body).Decode(&repository); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &repository, nil
}
