package models

import "github.com/google/go-github/github"

type PullRequestStatus struct {
	*github.CombinedStatus
	Mergeable bool
	MergeableStatus string
}