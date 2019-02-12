package interfaces

import "github.com/google/go-github/github"

type PullRequestAPI interface {
	GetPullRequest(pr int, repo string) (*github.PullRequest, error)
}
