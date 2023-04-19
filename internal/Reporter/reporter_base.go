package Reporter

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/google/go-github/v50/github"
	"golang.org/x/oauth2"
)

//nolint:revive
type ReporterType = string

const (
	GithubReporterType ReporterType = "github"
	GenericReporterType    ReporterType = "generic"
)

type Reporter interface {
	PostReport() error
}

type BasicReporter struct {
	Report string
}

func (r *BasicReporter) PostReport() error {
	fmt.Println(r.Report)
	return nil
}

func NewReporter(reporterType string, report string) (Reporter, error) {

	switch reporterType {
	case GithubReporterType:
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
	case GenericReporterType:
		// Generic Reporter will just post to stdout always.
		return &BasicReporter{
			report,
		}, nil
	}
	return nil, errors.New("Incompatible Reporter specific. Must be github or basic")
}
