package PlanAnalyzer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEmojisEmptyReturn(t *testing.T) {
	input := map[string][]string{
		"Create": {},
		"Update": {},
	}
	result := getEmojis(input)
	assert.Equal(t, result, "", "Result should be empty because changeset is empty")
}

func TestGetEmojisNotEmptyReturn(t *testing.T) {
	input := map[string][]string{
		"Create": {"one", "two", "three"},
		"Update": {},
	}
	result := getEmojis(input)
	assert.Equal(t, result, ":pencil2:", "Result should contain pencil")
}

func TestGetGitDiffMatch(t *testing.T) {
	result, _ := getGitDiff("Create")

	assert.Equal(t, result, "+", "The result of input Create should return +")
}

func TestGetGitDiffNoMatch(t *testing.T) {
	_, exists := getGitDiff("TestNotAMatch")

	assert.Equal(t, exists, false, "The exists boolean should be false because key is not a match")
}

func TestGetSortedWorkspaces(t *testing.T) {
	var changeSet = map[string][]string{
		"create": {"resource1"},
	}

	planAnalyzer := &PlanAnalyzer{
		[]PlanExtended{},
		[][]string{{"Workspace", "To Create", "To Update", "To Destroy", "To Replace"}},
		map[string]map[string][]string{
			"workspace4": {},
			"workspace1": {},
			"workspace3": {},
			"aworkspace": {},
		},
		changeSet,
	}

	var expected = []string{"aworkspace", "workspace1", "workspace3", "workspace4"}
	result := getSortedWorkspaces(planAnalyzer.Changes)
	assert.Equal(t, expected, result)
}
