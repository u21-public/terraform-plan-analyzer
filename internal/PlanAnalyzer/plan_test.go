package PlanAnalyzer

import (
	"testing"
	"fmt"
	"errors"
	"os"
	"path/filepath"
	"github.com/stretchr/testify/assert"
)

func TestParseWorkspaceNameSuccess(t *testing.T) {
	workspace, _ := ParseWorkspaceName("tfplan-account-region-environment.json")
    assert.Equal(t, workspace, "account-region-environment", "Result should be account-region-environment")
}

func TestParseWorkspaceNameEmptyStringErr(t *testing.T) {
	_, err := ParseWorkspaceName("")
    assert.Equal(t, err, errors.New("filename given was empty string"), "Result should error out with: filename given was empty string")
}

func TestFilePathWalkDirSuccess(t *testing.T) {
	var expected_files_list []string
	filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			expected_files_list = append(expected_files_list, path)
		}
		return nil
	})

	result, _ := FilePathWalkDir(".")
	assert.Equal(t, result, expected_files_list, "Result should contain all file names in current directory")
}

func TestReadPlansSuccess(t *testing.T) {
	result := ReadPlans("test_data/tfplan-account-region-environment.json")
	fmt.Println(result[0])
	fmt.Println([]PlanExtended(result))
}