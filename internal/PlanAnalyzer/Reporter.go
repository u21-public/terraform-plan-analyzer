
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
	Repo string
	Owner string
}

type Reporter struct {
	githubEnabled  bool
	GithubReporter *GithubReporter
	Report 		   string
}

func (r *Reporter) GetReportComment(issue int) *github.IssueComment {
	ctx := context.Background()

	opt := &github.IssueListCommentsOptions{
		ListOptions: github.ListOptions{PerPage: 30},
	}

	for {
		comments, resp, err := r.GithubReporter.Client.Issues.ListComments(ctx,r.GithubReporter.Owner , r.GithubReporter.Repo,issue, opt)
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

func (r *Reporter) PostReport(){
	if(r.githubEnabled) {
		ctx := context.Background()
		issue, err := strconv.Atoi(r.GithubReporter.Issue)
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
			_, _, err = r.GithubReporter.Client.Issues.EditComment(ctx,r.GithubReporter.Owner , r.GithubReporter.Repo, id, &report)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
		} else {
			_, _, err = r.GithubReporter.Client.Issues.CreateComment(ctx,r.GithubReporter.Owner , r.GithubReporter.Repo,issue, &report)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
		}

		fmt.Println("Posted to Github!, (Repo:", r.GithubReporter.Repo, ", Issue:", r.GithubReporter.Issue, ")")

	}else {
		fmt.Println(r.Report)
	}
}

func NewReporter(githubEnabled bool, report string) (*Reporter, error) {

	githubToken, gtPresent := os.LookupEnv("GITHUB_TOKEN")
	githubOwner, goPresent := os.LookupEnv("GITHUB_OWNER")
	githubRepo, grPresent := os.LookupEnv("GITHUB_REPOSITORY")
	githubIssue, gnPresent := os.LookupEnv("GITHUB_PR_NUMBER")

	if(githubEnabled){
		if (!gtPresent){
			return nil, errors.New("Error: GITHUB_TOKEN not set. Can't initialize Github Integration! Set ENVs or disable github integration.")
		}
		if (!grPresent){
			return nil, errors.New("Error: GITHUB_REPOSITORY not set. Can't initialize Github Integration! Set ENVs or disable github integration.")
		}
		if (!gnPresent){
			return nil, errors.New("Error: GITHUB_PR_NUMBER not set. Can't initialize Github Integration! Set ENVs or disable github integration. ")
		}
		if (!goPresent) {
			return nil, errors.New("Error: GITHUB_OWNER not set. Can't initialize Github Integration! Set ENVs or disable github integration. ")
		}
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
	}

	return &Reporter{
		githubEnabled,
		githubReporter,
		report,
	}, nil
}
