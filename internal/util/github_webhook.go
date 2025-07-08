package util

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
)

type GitHubWebhookUtil struct {
	Secret string
	Log    *logrus.Logger
}

func NewGitHubWebhookUtil(secret string, log *logrus.Logger) *GitHubWebhookUtil {
	return &GitHubWebhookUtil{
		Secret: secret,
		Log:    log,
	}
}

// VerifyWebhookSignature verifies the GitHub webhook signature
func (g *GitHubWebhookUtil) VerifyWebhookSignature(r *http.Request, body []byte) bool {
	signature := r.Header.Get("X-Hub-Signature-256")
	if signature == "" {
		g.Log.Warn("No GitHub signature found in request")
		return false
	}

	// Remove "sha256=" prefix
	signature = strings.TrimPrefix(signature, "sha256=")

	// Create HMAC
	h := hmac.New(sha256.New, []byte(g.Secret))
	h.Write(body)
	expectedSignature := hex.EncodeToString(h.Sum(nil))

	return hmac.Equal([]byte(signature), []byte(expectedSignature))
}

// GetWebhookURL returns the webhook URL for a given server
func (g *GitHubWebhookUtil) GetWebhookURL(serverURL string) string {
	return fmt.Sprintf("%s/webhook/github", serverURL)
}

// GetWebhookPayloadURL returns the webhook payload URL for GitHub API
func (g *GitHubWebhookUtil) GetWebhookPayloadURL(owner, repo string) string {
	return fmt.Sprintf("https://api.github.com/repos/%s/%s/hooks", owner, repo)
}
