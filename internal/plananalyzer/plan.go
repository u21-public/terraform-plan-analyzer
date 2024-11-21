package plananalyzer

import (
	"encoding/json"
	"errors"
	"io"
	"log"
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
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

func ParseWorkspaceName(planFileName string, formatByFolder bool) (string, error) {
	var workspaceName string
	var planBaseName string

	if formatByFolder {
		planBaseName = planFileName
	} else {
		planBaseName = filepath.Base(planFileName)
	}

	if planBaseName == "." {
		return "", errors.New("filename given was empty string")
	}

	planBaseNameSplit := strings.Split(planBaseName, ".json")
	if len(planBaseName) == 1 {
		return "", errors.New("plan filename must have a .json extension")
	}

	planNoExt := planBaseNameSplit[0]

	if formatByFolder {
		planNoPrefixSplit := strings.Split(planNoExt, "/tfplan")
		if len(planNoPrefixSplit) == 1 {
			return "", errors.New("plan filename must be tfplan.json")
		}

		workspaceName = planNoPrefixSplit[0]
	} else {
		planNoPrefixSplit := strings.Split(planNoExt, "tfplan-")

		if len(planNoPrefixSplit) > 1 {
			workspaceName = planNoPrefixSplit[1]
			_ = workspaceName
		} else {
			return "", errors.New("plan filename must be prefixed with tfplan-")
		}
	}

	return workspaceName, nil
}

func ReadPlans(plansFolderPath string) []PlanExtended {
	var plans []PlanExtended

	log.Println("Reading the plans in...`", plansFolderPath, "`")
	files, err := FilePathWalkDir(plansFolderPath)

	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		log.Println(file)
		plan := PlanExtended{}
		jsonFile, err := os.Open(file)
		if err != nil {
			log.Println(err)
		}
		defer jsonFile.Close()
		byteValue, err := io.ReadAll(jsonFile)
		if err != nil {
			log.Println(err)
		}
		err = json.Unmarshal(byteValue, &plan)
		if err != nil {
			log.Println(err)
		}
		plan.Analyze()

		fileRelativePath, err := filepath.Rel(plansFolderPath, file)
		if err != nil {
			log.Println(err)
		}

		workspace, err := ParseWorkspaceName(fileRelativePath, true)
		if err != nil {
			log.Println(err, "Arguments given: ", file)
		}
		plan.Workspace = workspace
		plans = append(plans, plan)
	}

	return plans
}
