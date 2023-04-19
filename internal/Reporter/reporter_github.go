package Reporter

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/go-github/v50/github"
)

const (
	GithubAPIPageLimit int = 30
)

type GithubReporter struct {
	Client *github.Client
	Issue  string
	Repo   string
	Owner  string
	Report string
}

// Determines if a report already exists on a given PR, this will determine if we edit
// an existing comment or create a new comment
func (r *GithubReporter) GetReportComment(issue int) (*github.IssueComment, bool, error) {
	ctx := context.Background()

	opt := &github.IssueListCommentsOptions{
		ListOptions: github.ListOptions{PerPage: GithubAPIPageLimit},
	}

	for {
		comments, resp, err := r.Client.Issues.ListComments(ctx, r.Owner, r.Repo, issue, opt)
		if err != nil {
			return nil, false, err
		}

		for _, comment := range comments {
			// We use the title of the report as an "id" to match against.
			// This breaks if multiple comments have this string, and it will only return
			// one of those comments.
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
	} else {
		_, _, err = r.Client.Issues.CreateComment(ctx, r.Owner, r.Repo, issue, &report)
	}
	if err != nil {
		return err
	}

	fmt.Println("Posted to Github!, (Repo:", r.Repo, ", Issue:", r.Issue, ")")
	return nil
}
