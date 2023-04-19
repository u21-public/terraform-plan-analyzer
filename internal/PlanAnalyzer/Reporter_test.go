package PlanAnalyzer

import (
	"errors"
	"github.com/google/go-github/v50/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"os"
	"testing"
)

func TestNewReporterNoGhToken(t *testing.T) {
	os.Unsetenv("GITHUB_TOKEN")
	os.Setenv("GITHUB_REPOSITORY", "foobar")
	os.Setenv("GITHUB_PR_NUMBER", "foobar")
	os.Setenv("GITHUB_OWNER", "foobar")

	_, err := NewReporter("github", "")
	expectedResult := errors.New("error: GITHUB_TOKEN not set. Can't initialize Github Integration! Set ENVs or disable github integration")
	assert.Equal(t, err, expectedResult, "Result should contain pencil")
}

func TestNewReporterNoGhOwner(t *testing.T) {
	os.Unsetenv("GITHUB_OWNER")
	os.Setenv("GITHUB_REPOSITORY", "foobar")
	os.Setenv("GITHUB_PR_NUMBER", "foobar")
	os.Setenv("GITHUB_TOKEN", "foobar")

	_, err := NewReporter("github", "")
	expectedResult := errors.New("error: GITHUB_OWNER not set. Can't initialize Github Integration! Set ENVs or disable github integration")
	assert.Equal(t, err, expectedResult, "Result should contain pencil")
}

func TestNewReporterNoGhRepo(t *testing.T) {
	os.Unsetenv("GITHUB_REPOSITORY")
	os.Setenv("GITHUB_TOKEN", "foobar")
	os.Setenv("GITHUB_PR_NUMBER", "foobar")
	os.Setenv("GITHUB_OWNER", "foobar")

	_, err := NewReporter("github", "")
	expectedResult := errors.New("error: GITHUB_REPOSITORY not set. Can't initialize Github Integration! Set ENVs or disable github integration")
	assert.Equal(t, err, expectedResult, "Result should contain pencil")
}

func TestNewReporterNoGhPrNumber(t *testing.T) {
	os.Unsetenv("GITHUB_PR_NUMBER")
	os.Setenv("GITHUB_REPOSITORY", "foobar")
	os.Setenv("GITHUB_TOKEN", "foobar")
	os.Setenv("GITHUB_OWNER", "foobar")

	_, err := NewReporter("github", "")
	expectedResult := errors.New("error: GITHUB_PR_NUMBER not set. Can't initialize Github Integration! Set ENVs or disable github integration")
	assert.Equal(t, err, expectedResult, "Result should contain pencil")
}

func TestGetReportCommentNoComments(t *testing.T) {
	mockedHTTPClient := mock.NewMockedHTTPClient(
		mock.WithRequestMatch(
			mock.GetReposIssuesCommentsByOwnerByRepoByIssueNumber,
			[]*github.IssueComment{},
		),
	)
	client := github.NewClient(mockedHTTPClient)
	githubReporter := &GithubReporter{
		client,
		"1",
		"terraform-plan-analyzer",
		"u21-public",
		"report",
	}

	_, found, _ := githubReporter.GetReportComment(1)
	assert.Equal(t, false, found, "Wrong")
}

func TestGetReportCommentListCommentError(t *testing.T) {
	failMessage := "GetCommentsRequestFailed"
	mockedHTTPClient := mock.NewMockedHTTPClient(
		mock.WithRequestMatchHandler(
			mock.GetReposIssuesCommentsByOwnerByRepoByIssueNumber,
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				mock.WriteError(
					w,
					http.StatusInternalServerError,
					failMessage,
				)
			}),
		),
	)
	client := github.NewClient(mockedHTTPClient)
	githubReporter := &GithubReporter{
		client,
		"1",
		"terraform-plan-analyzer",
		"u21-public",
		"report",
	}

	_, _, err := githubReporter.GetReportComment(1)
	assert.Equal(t, err.(*github.ErrorResponse).Message, failMessage)
}

func TestGetReportCommentListOnePage(t *testing.T) {
	mockedHTTPClient := mock.NewMockedHTTPClient(
		mock.WithRequestMatch(
			mock.GetReposIssuesCommentsByOwnerByRepoByIssueNumber,
			[]*github.IssueComment{
				{
					Body: github.String("Terraform Plan Analyzer Report"),
				},
			},
		),
	)
	client := github.NewClient(mockedHTTPClient)
	githubReporter := &GithubReporter{
		client,
		"1",
		"terraform-plan-analyzer",
		"u21-public",
		"report",
	}

	issue, found, _ := githubReporter.GetReportComment(1)
	assert.Equal(t, true, found, "Wrong")
	assert.Equal(t, "Terraform Plan Analyzer Report", *issue.Body)
}

func TestGetReportCommentListTwoPages(t *testing.T) {
	var mockIssueComments = []*github.IssueComment{}
	var mockComment = github.IssueComment{
		Body: github.String("foobar"),
	}

	for i := 0; i < 30; i++ {
		mockIssueComments = append(mockIssueComments, &mockComment)
	}
	mockIssueComments = append(mockIssueComments, &github.IssueComment{
		Body: github.String("Terraform Plan Analyzer Report"),
		ID:   github.Int64(123),
	})

	mockedHTTPClient := mock.NewMockedHTTPClient(
		mock.WithRequestMatch(
			mock.GetReposIssuesCommentsByOwnerByRepoByIssueNumber,
			mockIssueComments,
		),
	)
	client := github.NewClient(mockedHTTPClient)
	githubReporter := &GithubReporter{
		client,
		"1",
		"terraform-plan-analyzer",
		"u21-public",
		"report",
	}

	issue, found, _ := githubReporter.GetReportComment(1)
	assert.Equal(t, true, found, "Wrong")
	assert.Equal(t, "Terraform Plan Analyzer Report", *issue.Body)
	assert.Equal(t, *github.Int64(123), *issue.ID)

}

func TestPostReportIssueIntConvertError(t *testing.T) {
	mockedHTTPClient := mock.NewMockedHTTPClient(
		mock.WithRequestMatch(
			mock.GetReposIssuesCommentsByOwnerByRepoByIssueNumber,
			[]*github.IssueComment{
				{
					Body: github.String("Terraform Plan Analyzer Report"),
				},
			},
		),
	)
	client := github.NewClient(mockedHTTPClient)

	githubReporter := &GithubReporter{
		client,
		"asd",
		"asd",
		"asd",
		"asd",
	}

	err := githubReporter.PostReport()
	assert.EqualErrorf(t, err, "strconv.Atoi: parsing \"asd\": invalid syntax", "")

}

func TestPostReportIssueEditError(t *testing.T) {
	failMessage := "EditCommentsRequestFailed"
	mockedHTTPClient := mock.NewMockedHTTPClient(
		mock.WithRequestMatch(
			mock.GetReposIssuesCommentsByOwnerByRepoByIssueNumber,
			[]*github.IssueComment{
				{
					Body: github.String("Terraform Plan Analyzer Report"),
				},
			},
		),
		mock.WithRequestMatchHandler(
			mock.PatchReposIssuesCommentsByOwnerByRepoByCommentId,
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				mock.WriteError(
					w,
					http.StatusInternalServerError,
					failMessage,
				)
			}),
		),
	)
	client := github.NewClient(mockedHTTPClient)
	githubReporter := &GithubReporter{
		client,
		"1",
		"terraform-plan-analyzer",
		"u21-public",
		"report",
	}

	err := githubReporter.PostReport()
	assert.Equal(t, err.(*github.ErrorResponse).Message, failMessage, "")
}

func TestPostReportIssueCreateError(t *testing.T) {
	failMessage := "CreateCommentsRequestFailed"
	mockedHTTPClient := mock.NewMockedHTTPClient(
		mock.WithRequestMatch(
			mock.GetReposIssuesCommentsByOwnerByRepoByIssueNumber,
			[]*github.IssueComment{},
		),
		mock.WithRequestMatchHandler(
			mock.PostReposIssuesCommentsByOwnerByRepoByIssueNumber,
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				mock.WriteError(
					w,
					http.StatusInternalServerError,
					failMessage,
				)
			}),
		),
	)
	client := github.NewClient(mockedHTTPClient)
	githubReporter := &GithubReporter{
		client,
		"1",
		"terraform-plan-analyzer",
		"u21-public",
		"report",
	}

	err := githubReporter.PostReport()
	assert.Equal(t, err.(*github.ErrorResponse).Message, failMessage, "")
}

func TestPostReportIssueGetReportCommentError(t *testing.T) {
	failMessage := "GetCommentsRequestFailed"
	mockedHTTPClient := mock.NewMockedHTTPClient(
		mock.WithRequestMatchHandler(
			mock.GetReposIssuesCommentsByOwnerByRepoByIssueNumber,
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				mock.WriteError(
					w,
					http.StatusInternalServerError,
					failMessage,
				)
			}),
		),
	)
	client := github.NewClient(mockedHTTPClient)
	githubReporter := &GithubReporter{
		client,
		"1",
		"terraform-plan-analyzer",
		"u21-public",
		"report",
	}
	err := githubReporter.PostReport()
	assert.Equal(t, err.(*github.ErrorResponse).Message, failMessage, "")
}

func TestPostReportIssue(t *testing.T) {
	mockedHTTPClient := mock.NewMockedHTTPClient(
		mock.WithRequestMatch(
			mock.GetReposIssuesCommentsByOwnerByRepoByIssueNumber,
			[]*github.IssueComment{
				{
					Body: github.String("Terraform Plan Analyzer Report"),
				},
			},
		),
		mock.WithRequestMatch(
			mock.PatchReposIssuesCommentsByOwnerByRepoByCommentId,
			github.IssueComment{},
			nil,
			nil,
		),
	)

	client := github.NewClient(mockedHTTPClient)
	githubReporter := &GithubReporter{
		client,
		"1",
		"terraform-plan-analyzer",
		"u21-public",
		"report",
	}
	err := githubReporter.PostReport()
	assert.Equal(t, err, nil)
}
