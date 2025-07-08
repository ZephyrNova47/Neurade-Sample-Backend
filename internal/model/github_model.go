package model

type GitHubPullRequest struct {
	ID      int    `json:"id"`
	Number  int    `json:"number"`
	Title   string `json:"title"`
	Body    string `json:"body"`
	State   string `json:"state"`
	HTMLURL string `json:"html_url"`
	User    struct {
		Login string `json:"login"`
	} `json:"user"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type GitHubRepository struct {
	FullName string `json:"full_name"`
	HTMLURL  string `json:"html_url"`
}

type FetchPullRequestsRequest struct {
	CourseID    int    `json:"course_id"`
	GithubURL   string `json:"github_url"`
	GithubToken string `json:"github_token"`
}

type FetchPullRequestsResponse struct {
	Message           string `json:"message"`
	PullRequestsCount int    `json:"pull_requests_count"`
	CourseID          int    `json:"course_id"`
}
