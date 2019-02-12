package github

import (
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"os"
	"context"
)

func NewConnection() *github.Client{
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_ACCESS_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}
