package PlanAnalyzer

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	tfjson "github.com/hashicorp/terraform-json"
)

type PlanExtended struct {
	tfjson.Plan
	ChangeSet []*tfjson.ResourceChange
	ToUpdate  []string
	ToCreate  []string
	ToDestroy []string
	ToReplace []string
	Workspace string
}

func (p *PlanExtended) Analyze() {
	for _, change := range p.ResourceChanges {

		// Organize Changes into logical actions for quick
		// look up later
		if change.Change.Actions.Create() {
			p.ToCreate = append(p.ToCreate, change.Address)
		} else if change.Change.Actions.Delete() {
			p.ToDestroy = append(p.ToDestroy, change.Address)
		} else if change.Change.Actions.Update() {
			p.ToUpdate = append(p.ToUpdate, change.Address)
		} else if change.Change.Actions.Replace() {
			p.ToReplace = append(p.ToReplace, change.Address)
		}

		// Track all changes for quick look up later
		if !change.Change.Actions.NoOp() {
			p.ChangeSet = append(p.ChangeSet, change)
		}
	}
}

func (p *PlanExtended) getActions() map[string][]string {
	var changeSet = map[string][]string{
		Create:  p.ToCreate,
		Destroy: p.ToDestroy,
		Replace: p.ToReplace,
		Update:  p.ToUpdate,
	}
	return changeSet
}

func FilePathWalkDir(root string) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

func ParseWorkspaceName(planFileName string) (string, error) {
	planBaseName := filepath.Base(planFileName)

	if planBaseName == "." {
		return "", errors.New("filename given was empty string")
	}

	planNoExt := strings.Split(planBaseName, ".json")[0]
	planNoPrefix := strings.Split(planNoExt, "tfplan-")[1]

	if len(planNoPrefix) == 1 {
		return "", errors.New("plan filename must be prefixed with tfplan-")
	}

	return planNoPrefix, nil
}

func ReadPlans(plansFolderPath string) []PlanExtended {
	var plans []PlanExtended

	fmt.Println("Reading the plans in...`", plansFolderPath, "`")
	files, err := FilePathWalkDir(plansFolderPath)
	if err != nil {
		fmt.Println(err, "Arguments passed: ", plansFolderPath)
		os.Exit(1)
	}
	for _, file := range files {
		plan := PlanExtended{}
		jsonFile, err := os.Open(file)
		if err != nil {
			fmt.Println(err)
		}
		defer jsonFile.Close()
		byteValue, _ := io.ReadAll(jsonFile)
		err_two := plan.UnmarshalJSON([]byte(byteValue))
		if err_two != nil {
			fmt.Println(err_two)
		}
		plan.Analyze()
		workspace, err_parse_workspace := ParseWorkspaceName(file)
		if err_parse_workspace != nil {
			fmt.Println(err_parse_workspace, "Arguments given: ", file)
		}
		plan.Workspace = workspace
		plans = append(plans, plan)
	}

	return plans
}
