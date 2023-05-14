package PlanAnalyzer

import (
	"bytes"
	"errors"
	"log"
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
	assert.Equal(t, err, errors.New("plan filename must be prefixed with tfplan-"),
		"Result should error out with: plan filename must be prefixed with tfplan-")
}

func TestFilePathWalkDirSuccess(t *testing.T) {
	directory := t.TempDir()
	_, err1 := os.CreateTemp(directory, "file1.tf")
	_, err2 := os.CreateTemp(directory, "file2.tf")
	if err1 != nil || err2 != nil {
		log.Fatal("Error occurred while creating temporarily files in directory")
	}
	filesFound, _ := FilePathWalkDir(directory)
	assert.Equal(t, len(filesFound), 2, "Result should be equivalent to 2")
}

func TestFilePathWalkDirEmptyDir(t *testing.T) {
	directory := t.TempDir()
	filesFound, _ := FilePathWalkDir(directory)
	assert.Equal(t, len(filesFound), 0, "Result should be equivalent to 0")
}

func TestFilePathWalkInvalidDir(t *testing.T) {
	_, err := FilePathWalkDir("/random_invalid_path")
	assert.Contains(t, err.Error(), "no such file or directory")
}

func TestFileFoldersInsideFolders(t *testing.T) {
	directory, _ := os.MkdirTemp("", "first_dir")
	secondaryDir, _ := os.MkdirTemp(directory, "second_dir")
	_, err := os.CreateTemp(secondaryDir, "file1.tf")
	if err != nil {
		log.Fatal("Error occurred while creating temporarily files in directory")
	}
	filesFound, _ := FilePathWalkDir(directory)
	assert.Equal(t, len(filesFound), 1, "Result should be equivalent to 1")
}

func TestReadPlansEmptyDir(t *testing.T) {
	directory := t.TempDir()
	plansList := ReadPlans(directory)
	assert.Equal(t, len(plansList), 0, "Plans list should return 0")
}

func TestReadPlansSuccess(t *testing.T) {
	directory := t.TempDir()
	destinationFile, err := os.CreateTemp(directory, "tfplan-file1.tf")
	if err != nil {
		log.Fatal("Error occurred while creating temporarily files in directory")
	}
	absPath, _ := filepath.Abs("../../examples/plans_json/basic_example/tfplan-example1-only-creates.json")
	input, _ := os.ReadFile(absPath)
	writeErr := os.WriteFile(destinationFile.Name(), input, 0644)
	if writeErr != nil {
		log.Fatal("Error occurred while writing files to destination file")
	}

	plansList := ReadPlans(directory)
	toCreate := plansList[0].ToCreate

	assert.Equal(t, len(toCreate), 3, "Plan should create 3 S3 example buckets")
}

func TestReadPlansInvalidWorkspaceName(t *testing.T) {
	directory := t.TempDir()
	_, err := os.CreateTemp(directory, "account-us-west-2-prod1.json")
	if err != nil {
		log.Fatal("Error occurred while creating temporarily files in directory")
	}

	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stderr)
	}()
	ReadPlans(directory)

	assert.Contains(t, buf.String(), "plan filename must be prefixed with tfplan-")
}
