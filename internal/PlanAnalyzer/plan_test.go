package PlanAnalyzer

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseWorkspaceNameSuccess(t *testing.T) {
	workspace, _ := ParseWorkspaceName("tfplan-account-region-environment.json")
	assert.Equal(t, workspace, "account-region-environment", "Result should be account-region-environment")
}

func TestParseWorkspaceNameEmptyString(t *testing.T) {
	_, err := ParseWorkspaceName("")
	assert.Equal(t, err, errors.New("filename given was empty string"), "Result should error out with: filename given was empty string")
}

func TestParseWorkspaceNameNoPrefix(t *testing.T) {
	_, err := ParseWorkspaceName("account-us-west-2-prod1.json")
	assert.Equal(t, err, errors.New("plan filename must be prefixed with tfplan-"), "Result should error out with: plan filename must be prefixed with tfplan-")
}

func TestFilePathWalkDirSuccess(t *testing.T) {
	// Mock os.FileInfo
	var expectedFilesList []string
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			expectedFilesList = append(expectedFilesList, path)
		}
		return nil
	})
	if err != nil {
		fmt.Println("Error with iterating current directory")
	}

	result, _ := FilePathWalkDir(".")
	assert.Equal(t, result, expectedFilesList, "Result should contain all file names in current directory")
}

// How to test for error path here? maybe we need better error handling
// func TestFilePathWalkDirError(t *testing.T) {
// 	// expected_files_list := []string{}
//     // files, _ := ioutil.ReadDir(".")
//     // for _, file := range files {
//     //     expected_files_list= append(expected_files_list, file.Name())
//     // }

// 	result, _ := FilePathWalkDir("/")
// 	fmt.Println(result)
// 	// assert.Equal(t, result, expected_files_list, "Result should contain all file names in current directory")
// }

// func TestReadPlansSuccess(t *testing.T) {
// 	result := ReadPlans("test_data/tfplan-account-region-environment.json")
// 	fmt.Println(result[0])
// 	fmt.Println([]PlanExtended(result))
// }
