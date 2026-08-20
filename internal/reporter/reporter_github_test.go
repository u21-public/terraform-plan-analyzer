package reporter

import (
	"net/http"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/google/go-github/v50/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/stretchr/testify/assert"
)

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

	reporter, _ := githubReporter.GetReportComment(1)
	assert.Nilf(t, reporter, "reporter should be nil")
}

func TestGetReportCommentListCommentError(t *testing.T) {
	failMessage := "GetCommentsRequestFailed"
	mockedHTTPClient := mock.NewMockedHTTPClient(
		mock.WithRequestMatchHandler(
			mock.GetReposIssuesCommentsByOwnerByRepoByIssueNumber,
			http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
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

	_, err := githubReporter.GetReportComment(1)
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

	issue, _ := githubReporter.GetReportComment(1)
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

	issue, _ := githubReporter.GetReportComment(1)
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
			http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
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
			http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
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
			http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
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

func TestTruncateReportLeavesSmallReportUnchanged(t *testing.T) {
	report := "Terraform Plan Analyzer Report"

	assert.Equal(t, report, truncateReport(report))
}

func TestTruncateReportFitsGithubLimit(t *testing.T) {
	report := strings.Repeat("a", githubCommentMaxBytes+1)

	truncated := truncateReport(report)

	assert.Len(t, truncated, githubCommentMaxBytes)
	assert.True(t, strings.HasSuffix(truncated, truncatedReportNotice))
}

func TestTruncateReportPreservesUnicodeBoundaries(t *testing.T) {
	report := strings.Repeat("🙂", githubCommentMaxBytes)

	truncated := truncateReport(report)

	assert.LessOrEqual(t, len(truncated), githubCommentMaxBytes)
	assert.True(t, utf8.ValidString(truncated))
	assert.True(t, strings.HasSuffix(truncated, truncatedReportNotice))
}
