package interfaces

import (
	"github.com/google/go-github/github"
	"github.com/vincenciusgeraldo/sibyl/pkg/models"
)

type PullRequestAPI interface {
	GetPullRequest(pr int, repo string) (*github.PullRequest, error)
	GetPullRequestStatus(pr int, repo string) (models.PullRequestStatus, error)
}
