package PlanAnalyzer

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/google/go-github/v50/github"
	"golang.org/x/oauth2"
)

type Reporter interface {
	PostReport() error
}

type GithubReporter struct {
	Client *github.Client
	Issue  string
	Repo   string
	Owner  string
	Report string
}

type BasicReporter struct {
	Report string
}

func (r *BasicReporter) PostReport() error {
	fmt.Println(r.Report)
	return nil
}

func (r *GithubReporter) GetReportComment(issue int) (*github.IssueComment, bool, error) {
	ctx := context.Background()

	opt := &github.IssueListCommentsOptions{
		ListOptions: github.ListOptions{PerPage: 30},
	}

	for {
		comments, resp, err := r.Client.Issues.ListComments(ctx, r.Owner, r.Repo, issue, opt)
		if err != nil {
			return nil, false, err
		}

		for _, comment := range comments {
			if strings.Contains(comment.GetBody(), "Terraform Plan Analyzer Report") {
				return comment, true, nil
			}
		}

		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	return nil, false, nil
}

func (r *GithubReporter) PostReport() error {
	ctx := context.Background()
	issue, err := strconv.Atoi(r.Issue)
	if err != nil {
		return err
	}

	report := github.IssueComment{
		Body: &r.Report,
	}

	reportComment, found, err := r.GetReportComment(issue)
	if err != nil && !found {
		return err
	}
	if reportComment != nil {
		id := reportComment.GetID()
		_, _, err = r.Client.Issues.EditComment(ctx, r.Owner, r.Repo, id, &report)
		if err != nil {
			return err
		}
	} else {
		_, _, err = r.Client.Issues.CreateComment(ctx, r.Owner, r.Repo, issue, &report)
		if err != nil {
			return err
		}
	}

	fmt.Println("Posted to Github!, (Repo:", r.Repo, ", Issue:", r.Issue, ")")
	return nil
}

func NewReporter(reporterType string, report string) (Reporter, error) {

	switch reporterType {
	case "github":
		githubToken, gtPresent := os.LookupEnv("GITHUB_TOKEN")
		githubOwner, goPresent := os.LookupEnv("GITHUB_OWNER")
		githubRepo, grPresent := os.LookupEnv("GITHUB_REPOSITORY")
		githubIssue, gnPresent := os.LookupEnv("GITHUB_PR_NUMBER")

		if !gtPresent {
			return nil, errors.New("error: GITHUB_TOKEN not set. Can't initialize Github Integration! Set ENVs or disable github integration")
		}
		if !grPresent {
			return nil, errors.New("error: GITHUB_REPOSITORY not set. Can't initialize Github Integration! Set ENVs or disable github integration")
		}
		if !gnPresent {
			return nil, errors.New("error: GITHUB_PR_NUMBER not set. Can't initialize Github Integration! Set ENVs or disable github integration")
		}
		if !goPresent {
			return nil, errors.New("error: GITHUB_OWNER not set. Can't initialize Github Integration! Set ENVs or disable github integration")
		}

		ctx := context.Background()

		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: githubToken},
		)
		tc := oauth2.NewClient(ctx, ts)
		client := github.NewClient(tc)

		githubReporter := &GithubReporter{
			client,
			githubIssue,
			githubRepo,
			githubOwner,
			report,
		}
		return githubReporter, nil
	case "basic":
		return &BasicReporter{
			report,
		}, nil
	}
	return nil, errors.New("Imcompatable Reporter specific. Must be github or basic")
}
