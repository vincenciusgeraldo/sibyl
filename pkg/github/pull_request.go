package github

import (
	"github.com/google/go-github/github"
	"os"
	"context"
	"github.com/vincenciusgeraldo/sibyl/pkg/models"
)

type PullRequestAPI struct {
	cl    *github.Client
	owner string
}

func NewPullRequestInstance(cl *github.Client) *PullRequestAPI {
	return &PullRequestAPI{cl: cl, owner: os.Getenv("REPO_OWNER")}
}

func (pra *PullRequestAPI)GetPullRequest(pr int, repo string) (*github.PullRequest, error) {
	pullRequest, _, err := pra.cl.PullRequests.Get(context.Background(), pra.owner, repo, pr)
	if err != nil {
		return &github.PullRequest{}, err
	}

	return pullRequest, nil
}

func (pra *PullRequestAPI)GetPullRequestStatus(pr int, repo string) (models.PullRequestStatus, error) {
	pullRequest, err := pra.GetPullRequest(pr, repo)
	if err != nil {
		return models.PullRequestStatus{}, err
	}

	checkStatuses, _, err := pra.cl.Repositories.GetCombinedStatus(context.Background(), pra.owner, repo, *pullRequest.Head.SHA, &github.ListOptions{})
	if err != nil {
		return models.PullRequestStatus{}, err
	}

	return models.PullRequestStatus{
		CombinedStatus: checkStatuses,
		Mergeable: pullRequest.GetMergeable(),
		MergeableStatus: pullRequest.GetMergeableState(),
	}, nil
}