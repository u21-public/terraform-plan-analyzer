
package PlanAnalyzer

import "github.com/google/go-github/v50/github"
import "golang.org/x/oauth2"
import "context"
import "os"
import "errors"
import "fmt"
import 	"strconv"
import "strings"


type GithubReporter struct {
	Client  *github.Client
	Issue   string
	Repo    string
	Owner   string
	Report  string
}


func (r *GithubReporter) GetReportComment(issue int) *github.IssueComment {
	ctx := context.Background()

	opt := &github.IssueListCommentsOptions{
		ListOptions: github.ListOptions{PerPage: 30},
	}

	for {
		comments, resp, err := r.Client.Issues.ListComments(ctx,r.Owner , r.Repo,issue, opt)
		if err != nil {
			os.Exit(1)
		}

		for _,comment := range(comments) {
			if strings.Contains(comment.GetBody(), "Terraform Plan Analyzer Report") {
				return comment
			}
		}

		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	return nil
}

func (r *GithubReporter) PostReport(){
	ctx := context.Background()
	issue, err := strconv.Atoi(r.Issue)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	report := github.IssueComment{
		Body: &r.Report,
	}

	reportComment := r.GetReportComment(issue)
	if reportComment != nil {
		id := reportComment.GetID()
		_, _, err = r.Client.Issues.EditComment(ctx,r.Owner , r.Repo, id, &report)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	} else {
		_, _, err = r.Client.Issues.CreateComment(ctx,r.Owner , r.Repo,issue, &report)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	}

	fmt.Println("Posted to Github!, (Repo:", r.Repo, ", Issue:", r.Issue, ")")
}


func NewGithubReporter(report string) (*GithubReporter, error) {

	githubToken, gtPresent := os.LookupEnv("GITHUB_TOKEN")
	githubOwner, goPresent := os.LookupEnv("GITHUB_OWNER")
	githubRepo, grPresent := os.LookupEnv("GITHUB_REPOSITORY")
	githubIssue, gnPresent := os.LookupEnv("GITHUB_PR_NUMBER")

	if (!gtPresent){
		return nil, errors.New("error: GITHUB_TOKEN not set. Can't initialize Github Integration! Set ENVs or disable github integration.")
	}
	if (!grPresent){
		return nil, errors.New("error: GITHUB_REPOSITORY not set. Can't initialize Github Integration! Set ENVs or disable github integration.")
	}
	if (!gnPresent){
		return nil, errors.New("error: GITHUB_PR_NUMBER not set. Can't initialize Github Integration! Set ENVs or disable github integration.")
	}
	if (!goPresent) {
		return nil, errors.New("error: GITHUB_OWNER not set. Can't initialize Github Integration! Set ENVs or disable github integration.")
	}

	ctx := context.Background()

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: githubToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	
	githubReporter := &GithubReporter {
		client,
		githubIssue,
		githubRepo,
		githubOwner,
		report,
	}

	return githubReporter, nil
}
