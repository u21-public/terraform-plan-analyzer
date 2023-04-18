package PlanAnalyzer

import (
	"os"
	"testing"
	"errors"
	"github.com/stretchr/testify/assert"
)

func TestNewReporterNoGhToken(t *testing.T) {
	os.Unsetenv("GITHUB_TOKEN")
	os.Setenv("GITHUB_REPOSITORY", "foobar")
	os.Setenv("GITHUB_PR_NUMBER", "foobar")
	os.Setenv("GITHUB_OWNER", "foobar")

	_, err := NewGithubReporter("")
	expectedResult := errors.New("error: GITHUB_TOKEN not set. Can't initialize Github Integration! Set ENVs or disable github integration.")
	assert.Equal(t, err, expectedResult, "Result should contain pencil")
}

func TestNewReporterNoGhOwner(t *testing.T) {
	os.Unsetenv("GITHUB_OWNER")
	os.Setenv("GITHUB_REPOSITORY", "foobar")
	os.Setenv("GITHUB_PR_NUMBER", "foobar")
	os.Setenv("GITHUB_TOKEN", "foobar")

	_, err := NewGithubReporter("")
	expectedResult := errors.New("error: GITHUB_OWNER not set. Can't initialize Github Integration! Set ENVs or disable github integration.")
	assert.Equal(t, err, expectedResult, "Result should contain pencil")
}

func TestNewReporterNoGhRepo(t *testing.T) {
	os.Unsetenv("GITHUB_REPOSITORY")
	os.Setenv("GITHUB_TOKEN", "foobar")
	os.Setenv("GITHUB_PR_NUMBER", "foobar")
	os.Setenv("GITHUB_OWNER", "foobar")

	_, err := NewGithubReporter("")
	expectedResult := errors.New("error: GITHUB_REPOSITORY not set. Can't initialize Github Integration! Set ENVs or disable github integration.")
	assert.Equal(t, err, expectedResult, "Result should contain pencil")
}

func TestNewReporterNoGhPrNumber(t *testing.T) {
	os.Unsetenv("GITHUB_PR_NUMBER")
	os.Setenv("GITHUB_REPOSITORY", "foobar")
	os.Setenv("GITHUB_TOKEN", "foobar")
	os.Setenv("GITHUB_OWNER", "foobar")

	_, err := NewGithubReporter("")
	expectedResult := errors.New("error: GITHUB_PR_NUMBER not set. Can't initialize Github Integration! Set ENVs or disable github integration.")
	assert.Equal(t, err, expectedResult, "Result should contain pencil")
}
