package PlanAnalyzer

import (
	"fmt"
	"strconv"
	"time"
)

type PlanAnalyzer struct {
	Plans           []PlanExtended
	ComparisonTable [][]string

	// { <workspace>: {<action>: <resource>}}
	// ex: { prod1: { "ToUpdate": "resource1"}}
	UniqueChanges map[string]map[string][]string

	// { <action>: <resource>}
	// ex: { "ToUpdate: "resource1", "ToDestroy": "resource2" }
	SharedChanges map[string][]string
}

func (pa *PlanAnalyzer) ProcessPlans() {
	fmt.Println("Comparing Workspaces...")

	var hash = map[string]map[string]bool{
		Create:  {},
		Destroy: {},
		Update:  {},
		Replace: {},
	}

	var intersection = map[string][]string{
		Create:  {},
		Destroy: {},
		Update:  {},
		Replace: {},
	}

	// We run through all the plans and perform processing used for later
	// NOTE: we are doing multiple proccesses in same for loop for performance
	// reasons. Dont want to loop all changesets multiple times.
	for i, plan := range pa.Plans {
		pa.ComparisonTable = append(pa.ComparisonTable, []string{
			plan.Workspace,
			strconv.Itoa(len(plan.ToCreate)),
			strconv.Itoa(len(plan.ToUpdate)),
			strconv.Itoa(len(plan.ToDestroy)),
			strconv.Itoa(len(plan.ToReplace)),
		})
		pa.UniqueChanges[plan.Workspace] = plan.getActions()

		// Hash intersection for quick slice comparison
		for _, action := range SupportedAction {
			for _, address := range pa.UniqueChanges[plan.Workspace][action] {
				if i == 0 {
					hash[action][address] = true
				} else {
					if _, ok := hash[action][address]; ok {
						if i == len(pa.Plans)-1 {
							intersection[action] = append(intersection[action], address)
						}
					} else {
						delete(hash[action], address)
					}
				}
			}
		}
		pa.SharedChanges = intersection
	}
}

func (pa *PlanAnalyzer) GenerateLastUpdated() string {
	currentTime := time.Now()
	lastUpdated := fmt.Sprintf("Last Updated: `%s`\n\n", currentTime.String())
	return lastUpdated
}

func (pa *PlanAnalyzer) GenerateComparisonTable() string {
	var markdownTable string

	for row := range pa.ComparisonTable {
		for column := range pa.ComparisonTable[row] {
			markdownTable = markdownTable + "| "
			markdownTable = markdownTable + pa.ComparisonTable[row][column] + " "
		}
		markdownTable = markdownTable + "|"
		markdownTable = markdownTable + "\n"

		if row == 0 {
			markdownTable = markdownTable + "|-|:-:|:-:|:-:|:-:|\n"
		}
	}
	markdownTable = markdownTable + "\n\n"

	return markdownTable
}

func (pa *PlanAnalyzer) GenerateSharedResources() string {
	var sharedResources string

	sharedResources = sharedResources + "## All Workspaces" + getEmojis(pa.SharedChanges) + "\n"
	for action, changedResources := range pa.SharedChanges {

		result, _ := getGitDiff(action)
		// Open Code block
		sharedResources = sharedResources + "```diff\n"
		sharedResources = sharedResources + fmt.Sprintf("%s To %s %s\n", result, action, result)
		for _, resource := range changedResources {
			sharedResources = sharedResources + fmt.Sprintf("~ %s\n", resource)
		}
		// Close Code block
		sharedResources = sharedResources + "```\n\n"
	}

	return sharedResources
}

func (pa *PlanAnalyzer) GenerateUniqueResources() string {
	var UniqueChanges string

	UniqueChanges = UniqueChanges + "## Individual Workspaces\n"

	for workspace, changeSet := range pa.UniqueChanges {
		UniqueChanges = UniqueChanges + fmt.Sprintf("### %s %s\n", workspace, getEmojis(changeSet))
		for action, changedResources := range changeSet {

			result, _ := getGitDiff(action)
			// TODO: Do not show resources shared between between unique + shared resources (only unique changes)
			if len(changedResources) > 0 {
				UniqueChanges = UniqueChanges + "```diff\n"
				UniqueChanges = UniqueChanges + fmt.Sprintf("%s To %s %s\n", result, action, result)
				for _, resource := range changedResources {
					UniqueChanges = UniqueChanges + fmt.Sprintf("~ %s\n", resource)
				}
				UniqueChanges = UniqueChanges + "```\n\n"
			}
		}
	}

	return UniqueChanges
}

func (pa *PlanAnalyzer) GenerateReport() string {
	var report string

	reportTitle := fmt.Sprintf("# %s Terraform Plan Analyzer Report! %s\n", EmojiMap["title"], EmojiMap["title"])
	lastUpdated := pa.GenerateLastUpdated()
	markdownTable := pa.GenerateComparisonTable()
	sharedResources := pa.GenerateSharedResources()
	UniqueChanges := pa.GenerateUniqueResources()

	report = reportTitle + lastUpdated + markdownTable + sharedResources + UniqueChanges
	return report
}

func NewPlanAnalyzer(plans []PlanExtended) *PlanAnalyzer {
	return &PlanAnalyzer{
		plans,
		[][]string{{"Workspace", "To Create", "To Update", "To Destroy", "To Replace"}},
		map[string]map[string][]string{},
		map[string][]string{},
	}
}
