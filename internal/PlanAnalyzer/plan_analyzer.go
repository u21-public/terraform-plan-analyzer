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

	sharedResourceTitle := "## All Workspaces" + getEmojis(pa.SharedChanges) + "\n"
	sharedResources = sharedResourceTitle
	for _, action := range SupportedAction {

		changedResources := pa.SharedChanges[action]
		if len(changedResources) > 0 {
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
	}

	// Only occurs if no shared resources exist, so dont both just returning title
	if sharedResources == sharedResourceTitle {
		return ""
	}

	return sharedResources
}

func (pa *PlanAnalyzer) GenerateUniqueResources() string {
	var uniqueChanges string

	uniqueChangesTitle := "## Individual Workspaces\n"
	uniqueChanges = uniqueChanges + uniqueChangesTitle

	sortedWorkspaces := GetSortedWorkspaces(pa.UniqueChanges)
	for _, workspace := range sortedWorkspaces {
		uniqueChanges = uniqueChanges + pa.GenerateWorkspaceUniqueResources(workspace, pa.UniqueChanges[workspace])
	}

	if uniqueChanges == uniqueChangesTitle {
		return ""
	}
	return uniqueChanges
}

func (pa *PlanAnalyzer) GenerateWorkspaceUniqueResources(workspace string, changeSet map[string][]string) string {
	var uniqueChanges string
	resourceCount := 0

	uniqueChanges = uniqueChanges + fmt.Sprintf("### %s %s\n", workspace, getEmojis(changeSet))
	for _, action := range SupportedAction {
		changedResources := changeSet[action]
		result, _ := getGitDiff(action)

		if len(changedResources) > 0 {
			uniqueChanges = uniqueChanges + "```diff\n"
			uniqueChanges = uniqueChanges + fmt.Sprintf("%s To %s %s\n", result, action, result)
			for _, resource := range changedResources {
				if pa.IsChangeUnique(action, resource) {
					uniqueChanges = uniqueChanges + fmt.Sprintf("~ %s\n", resource)
					resourceCount = resourceCount + 1
				}
			}
			uniqueChanges = uniqueChanges + "```\n\n"
		}
	}

	if resourceCount == 0 {
		return ""
	}
	return uniqueChanges
}

func (pa *PlanAnalyzer) IsChangeUnique(action string, resource string) bool {
	for _, sharedResource := range pa.SharedChanges[action] {
		if resource == sharedResource {
			return false
		}
	}
	return true
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
