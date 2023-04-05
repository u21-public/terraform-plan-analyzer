package PlanAnalyzer

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestGetEmojisEmptyReturn(t *testing.T) {
	input := map[string][]string{
		"Create": []string{},
		"Update": []string{},
	}
	result := getEmojis(input)
	assert.Equal(t, result, "", "Result should be empty because changeset is empty")
}

func TestGetEmojisNotEmptyReturn(t *testing.T) {
	input := map[string][]string{
		"Create": []string{"one", "two", "three"},
		"Update": []string{"one"},
	}
	result := getEmojis(input)
	assert.Equal(t, result, ":pencil2::fountain_pen:", "Result should contain pencil and fountain pen")
}

func TestGetGitDiffMatch(t *testing.T) {
	result, _ := getGitDiff("Create")
	
	assert.Equal(t, result, "+", "The result of input Create should return +")
}

func TestGetGitDiffNoMatch(t *testing.T) {
	_ , exists := getGitDiff("TestNotAMatch")

	assert.Equal(t, exists, false, "The exists boolean should be false because key is not a match")
}