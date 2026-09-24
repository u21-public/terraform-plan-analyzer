package plananalyzer

import (
	"strings"
	"testing"

	tfjson "github.com/hashicorp/terraform-json"
	"github.com/stretchr/testify/assert"
)

func TestGenerateSharedResourcesEmpty(t *testing.T) {
	planAnalyzer := &PlanAnalyzer{
		[]PlanExtended{},
		[][]string{{"Workspace", "To Create", "To Update", "To Destroy", "To Replace"}},
		map[string]map[string][]string{},
		map[string][]string{},
	}

	result := planAnalyzer.GenerateSharedResources()
	assert.Equal(t, "", result)
}

func TestGenerateSharedResourcesNoResources(t *testing.T) {
	planAnalyzer := &PlanAnalyzer{
		[]PlanExtended{},
		[][]string{{"Workspace", "To Create", "To Update", "To Destroy", "To Replace"}},
		map[string]map[string][]string{},
		map[string][]string{Create: {}},
	}

	expected := ""
	result := planAnalyzer.GenerateSharedResources()
	assert.Equal(t, expected, result)
}

func TestGenerateSharedResourcesNotEmpty(t *testing.T) {
	planAnalyzer := &PlanAnalyzer{
		[]PlanExtended{},
		[][]string{{"Workspace", "To Create", "To Update", "To Destroy", "To Replace"}},
		map[string]map[string][]string{},
		map[string][]string{Create: {"resource1"}},
	}

	expected := "## All Workspaces:pencil2:\n```diff\n+ To Create +\n~ resource1\n```\n\n"
	result := planAnalyzer.GenerateSharedResources()
	assert.Equal(t, expected, result)
}

func TestGenerateResourcesNoUnique(t *testing.T) {
	var changeSet = map[string][]string{
		Create: {"resource1"},
	}

	planAnalyzer := &PlanAnalyzer{
		[]PlanExtended{},
		[][]string{{"Workspace", "To Create", "To Update", "To Destroy", "To Replace"}},
		map[string]map[string][]string{
			"workspace1": changeSet,
		},
		changeSet,
	}

	result := planAnalyzer.generateResources()
	assert.Equal(t, "", result)
}

func TestGenerateResourcesSomeUnique(t *testing.T) {
	var changeSet = map[string][]string{
		Create:  {"resource1"},
		Destroy: {"resource2"},
	}

	planAnalyzer := &PlanAnalyzer{
		[]PlanExtended{},
		[][]string{{"Workspace", "To Create", "To Update", "To Destroy", "To Replace"}},
		map[string]map[string][]string{
			"workspace1": changeSet,
		},
		map[string][]string{Create: {"resource1"}},
	}

	// Line is to long, so split it up
	expected := "## Individual Workspaces\n### workspace1 :pencil2::wastebasket:\n"
	expected = expected + "```diff\n- To Destroy -\n~ resource2\n```\n\n"
	result := planAnalyzer.generateResources()
	assert.Equal(t, expected, result)
}

func TestGenerateResourcesMultipleWorkspaces(t *testing.T) {
	var changeSet = map[string][]string{
		Create:  {"resource1"},
		Destroy: {"resource2"},
	}

	planAnalyzer := &PlanAnalyzer{
		[]PlanExtended{},
		[][]string{{"Workspace", "To Create", "To Update", "To Destroy", "To Replace"}},
		map[string]map[string][]string{
			"workspace1": changeSet,
			"workspace2": changeSet,
		},
		map[string][]string{Create: {"resource1"}},
	}

	// Line is to long, so split it up
	expected := "## Individual Workspaces\n### workspace1 :pencil2::wastebasket:\n"
	expected = expected + "```diff\n- To Destroy -\n~ resource2\n```\n\n"
	expected = expected + "### workspace2 :pencil2::wastebasket:\n```diff\n- To Destroy -\n~ resource2\n```\n\n"
	result := planAnalyzer.generateResources()
	assert.Equal(t, expected, result)
}

func TestGenerateWorkspaceResourcesMultipleWorkspacesSomeUnique(t *testing.T) {

	var changeSetOne = map[string][]string{
		Create: {"resource1"},
	}

	var changeSetTwo = map[string][]string{
		Create: {"resource1"},
		Update: {"resource2"},
	}

	planAnalyzer := &PlanAnalyzer{
		[]PlanExtended{},
		[][]string{{"Workspace", "To Create", "To Update", "To Destroy", "To Replace"}},
		map[string]map[string][]string{
			"workspace1": changeSetOne,
			"workspace2": changeSetTwo,
		},
		map[string][]string{Create: {"resource1"}},
	}

	expectedOne := ""
	expectedTwo := "### workspace2 :pencil2::fountain_pen:\n```diff\n! To Update !\n~ resource2\n```\n\n"

	resultOne := planAnalyzer.GenerateWorkspaceResources("workspace1", changeSetOne)
	resultTwo := planAnalyzer.GenerateWorkspaceResources("workspace2", changeSetTwo)

	assert.Equal(t, expectedOne, resultOne)
	assert.Equal(t, expectedTwo, resultTwo)
}

func TestGenerateWorkspaceResourcesMultipleWorkspacesNoUnique(t *testing.T) {

	var changeSetOne = map[string][]string{
		Create: {"resource1"},
	}

	var changeSetTwo = map[string][]string{
		Create: {"resource1"},
	}

	planAnalyzer := &PlanAnalyzer{
		[]PlanExtended{},
		[][]string{{"Workspace", "To Create", "To Update", "To Destroy", "To Replace"}},
		map[string]map[string][]string{
			"workspace1": changeSetOne,
			"workspace2": changeSetTwo,
		},
		map[string][]string{Create: {"resource1"}},
	}

	expectedOne := ""
	expectedTwo := ""

	resultOne := planAnalyzer.GenerateWorkspaceResources("workspace1", changeSetOne)
	resultTwo := planAnalyzer.GenerateWorkspaceResources("workspace2", changeSetTwo)

	assert.Equal(t, expectedOne, resultOne)
	assert.Equal(t, expectedTwo, resultTwo)
}

func TestGenerateWorkspaceResourcesMultipleWorkspacesOnlyUnique(t *testing.T) {

	var changeSetOne = map[string][]string{
		Create: {"resource1"},
	}

	var changeSetTwo = map[string][]string{
		Create: {"resource2"},
	}

	planAnalyzer := &PlanAnalyzer{
		[]PlanExtended{},
		[][]string{{"Workspace", "To Create", "To Update", "To Destroy", "To Replace"}},
		map[string]map[string][]string{
			"workspace1": changeSetOne,
			"workspace2": changeSetTwo,
		},
		map[string][]string{},
	}

	expectedOne := "### workspace1 :pencil2:\n```diff\n+ To Create +\n~ resource1\n```\n\n"
	expectedTwo := "### workspace2 :pencil2:\n```diff\n+ To Create +\n~ resource2\n```\n\n"

	resultOne := planAnalyzer.GenerateWorkspaceResources("workspace1", changeSetOne)
	resultTwo := planAnalyzer.GenerateWorkspaceResources("workspace2", changeSetTwo)

	assert.Equal(t, expectedOne, resultOne)
	assert.Equal(t, expectedTwo, resultTwo)
}

func TestGenerateWorkspaceResourcesNoUnique(t *testing.T) {

	var changeSet = map[string][]string{
		Create: {"resource1"},
	}

	planAnalyzer := &PlanAnalyzer{
		[]PlanExtended{},
		[][]string{{"Workspace", "To Create", "To Update", "To Destroy", "To Replace"}},
		map[string]map[string][]string{
			"workspace1": changeSet,
		},
		changeSet,
	}

	result := planAnalyzer.GenerateWorkspaceResources("workspace1", changeSet)
	assert.Equal(t, "", result)
}

func TestGenerateWorkspaceResourcesAllUnique(t *testing.T) {

	var changeSet = map[string][]string{
		Create:  {"resource1"},
		Destroy: {"resource2"},
	}

	planAnalyzer := &PlanAnalyzer{
		[]PlanExtended{},
		[][]string{{"Workspace", "To Create", "To Update", "To Destroy", "To Replace"}},
		map[string]map[string][]string{
			"workspace1": changeSet,
		},
		map[string][]string{Create: {}},
	}

	expected := "### workspace1 :pencil2::wastebasket:\n```diff\n+ To Create +\n~ resource1\n```\n\n```diff\n- To Destroy -\n~ resource2\n```\n\n"
	result := planAnalyzer.GenerateWorkspaceResources("workspace1", changeSet)
	assert.Equal(t, expected, result)
}

func TestGenerateWorkspaceResourcesSomeUnique(t *testing.T) {

	var changeSet = map[string][]string{
		Create:  {"resource1"},
		Destroy: {"resource2"},
	}

	planAnalyzer := &PlanAnalyzer{
		[]PlanExtended{},
		[][]string{{"Workspace", "To Create", "To Update", "To Destroy", "To Replace"}},
		map[string]map[string][]string{
			"workspace1": changeSet,
		},
		map[string][]string{Create: {"resource1"}},
	}

	expected := "### workspace1 :pencil2::wastebasket:\n```diff\n- To Destroy -\n~ resource2\n```\n\n"
	result := planAnalyzer.GenerateWorkspaceResources("workspace1", changeSet)
	assert.Equal(t, expected, result)
}

func TestIsChangeUnique(t *testing.T) {
	planAnalyzer := &PlanAnalyzer{
		[]PlanExtended{},
		[][]string{{"Workspace", "To Create", "To Update", "To Destroy", "To Replace"}},
		map[string]map[string][]string{},
		map[string][]string{
			Create: {"resource1"},
		},
	}
	result := planAnalyzer.IsChangeUnique(Create, "resource1")
	assert.Equal(t, false, result)
}

func TestIsChangeUniqueNotUnique(t *testing.T) {
	planAnalyzer := &PlanAnalyzer{
		[]PlanExtended{},
		[][]string{{"Workspace", "To Create", "To Update", "To Destroy", "To Replace"}},
		map[string]map[string][]string{},
		map[string][]string{
			Create: {"resource1"},
		},
	}
	result := planAnalyzer.IsChangeUnique(Create, "resource2")
	assert.Equal(t, true, result)
}

func TestIsChangeUniqueNotUniqueEmptyUnique(t *testing.T) {
	planAnalyzer := &PlanAnalyzer{
		[]PlanExtended{},
		[][]string{{"Workspace", "To Create", "To Update", "To Destroy", "To Replace"}},
		map[string]map[string][]string{},
		map[string][]string{},
	}
	result := planAnalyzer.IsChangeUnique(Create, "resource1")
	assert.Equal(t, true, result)
}

func TestProcessPlansSharedChangesNoShared(t *testing.T) {
	plans := []PlanExtended{
		{
			tfjson.Plan{},
			[]*tfjson.ResourceChange{},
			[]string{"resource1"},
			[]string{},
			[]string{},
			[]string{},
			"workspaceOne",
		},
		{
			tfjson.Plan{},
			[]*tfjson.ResourceChange{},
			[]string{"resource2"},
			[]string{},
			[]string{},
			[]string{},
			"workspaceTwo",
		},
	}

	planAnalyzer := &PlanAnalyzer{
		plans,
		[][]string{{"Workspace", "To Create", "To Update", "To Destroy", "To Replace"}},
		map[string]map[string][]string{},
		map[string][]string{},
	}
	planAnalyzer.ProcessPlans()
	result := map[string][]string(map[string][]string{"Create": {}, "Destroy": {}, "Replace": {}, "Update": {}})
	assert.Equal(t, result, planAnalyzer.SharedChanges)
}

func TestProcessPlansSharedChangesOneShared(t *testing.T) {
	plans := []PlanExtended{
		{
			tfjson.Plan{},
			[]*tfjson.ResourceChange{},
			[]string{"resource1"},
			[]string{},
			[]string{},
			[]string{},
			"workspaceOne",
		},
		{
			tfjson.Plan{},
			[]*tfjson.ResourceChange{},
			[]string{"resource1"},
			[]string{},
			[]string{},
			[]string{},
			"workspaceTwo",
		},
	}

	planAnalyzer := &PlanAnalyzer{
		plans,
		[][]string{{"Workspace", "To Create", "To Update", "To Destroy", "To Replace"}},
		map[string]map[string][]string{},
		map[string][]string{},
	}
	planAnalyzer.ProcessPlans()
	result := map[string][]string(map[string][]string{"Create": {}, "Destroy": {}, "Replace": {}, "Update": {"resource1"}})
	assert.Equal(t, result, planAnalyzer.SharedChanges)
}

func TestProcessPlansSharedChangesOneSharedAcrossMany(t *testing.T) {
	plans := []PlanExtended{
		{
			tfjson.Plan{},
			[]*tfjson.ResourceChange{},
			[]string{"resource1", "resource2", "resource3"},
			[]string{"resource2"},
			[]string{"resource1"},
			[]string{"resource1"},
			"workspaceOne",
		},
		{
			tfjson.Plan{},
			[]*tfjson.ResourceChange{},
			[]string{"resource1", "resource3"},
			[]string{},
			[]string{},
			[]string{},
			"workspaceTwo",
		},
		{
			tfjson.Plan{},
			[]*tfjson.ResourceChange{},
			[]string{"resource3"},
			[]string{},
			[]string{},
			[]string{},
			"workspaceTwo",
		},
		{
			tfjson.Plan{},
			[]*tfjson.ResourceChange{},
			[]string{"resource1", "resource3"},
			[]string{},
			[]string{},
			[]string{},
			"workspaceTwo",
		},
		{
			tfjson.Plan{},
			[]*tfjson.ResourceChange{},
			[]string{"resource2", "resource3"},
			[]string{},
			[]string{"resource1"},
			[]string{"resource4"},
			"workspaceTwo",
		},
	}

	planAnalyzer := &PlanAnalyzer{
		plans,
		[][]string{{"Workspace", "To Create", "To Update", "To Destroy", "To Replace"}},
		map[string]map[string][]string{},
		map[string][]string{},
	}
	planAnalyzer.ProcessPlans()
	result := map[string][]string{"Create": {}, "Destroy": {}, "Replace": {}, "Update": {"resource3"}}
	assert.Equal(t, result, planAnalyzer.SharedChanges)
}

func TestProcessPlansSharedChanges(t *testing.T) {
	plans := []PlanExtended{
		{
			tfjson.Plan{},
			[]*tfjson.ResourceChange{},
			[]string{"resource1", "resource2", "resource3", "resource4"},
			[]string{},
			[]string{},
			[]string{},
			"workspaceOne",
		},
		{
			tfjson.Plan{},
			[]*tfjson.ResourceChange{},
			[]string{"resource1"},
			[]string{},
			[]string{},
			[]string{},
			"workspaceTwo",
		},
	}

	planAnalyzer := &PlanAnalyzer{
		plans,
		[][]string{{"Workspace", "To Create", "To Update", "To Destroy", "To Replace"}},
		map[string]map[string][]string{},
		map[string][]string{},
	}
	planAnalyzer.ProcessPlans()
	result := map[string][]string{"Create": {}, "Destroy": {}, "Replace": {}, "Update": {"resource1"}}
	assert.Equal(t, result, planAnalyzer.SharedChanges)
}

func TestProcessPlansSharedPlansTwoSharedOneunique(t *testing.T) {
	plans := []PlanExtended{
		{
			tfjson.Plan{},
			[]*tfjson.ResourceChange{},
			[]string{"resource1"},
			[]string{},
			[]string{},
			[]string{},
			"workspaceOne",
		},
		{
			tfjson.Plan{},
			[]*tfjson.ResourceChange{},
			[]string{},
			[]string{},
			[]string{},
			[]string{},
			"workspaceTwo",
		},
		{
			tfjson.Plan{},
			[]*tfjson.ResourceChange{},
			[]string{"resource1"},
			[]string{},
			[]string{},
			[]string{},
			"workspaceTwo",
		},
	}

	planAnalyzer := &PlanAnalyzer{
		plans,
		[][]string{{"Workspace", "To Create", "To Update", "To Destroy", "To Replace"}},
		map[string]map[string][]string{},
		map[string][]string{},
	}
	planAnalyzer.ProcessPlans()
	result := map[string][]string{"Create": {}, "Destroy": {}, "Replace": {}, "Update": {}}
	assert.Equal(t, result, planAnalyzer.SharedChanges)
}

func TestGenerateComparisonTableEmpty(t *testing.T) {
	planAnalyzer := &PlanAnalyzer{
		[]PlanExtended{},
		[][]string{{"Workspace", "To Create", "To Update", "To Destroy", "To Replace"}},
		map[string]map[string][]string{},
		map[string][]string{},
	}
	expected := "| Workspace | To Create | To Update | To Destroy | To Replace |\n|-|:-:|:-:|:-:|:-:|\n\n\n"
	assert.Equal(t, expected, planAnalyzer.GenerateComparisonTable())
}

func TestGenerateComparisonTableWithWorkspaces(t *testing.T) {
	planAnalyzer := &PlanAnalyzer{
		[]PlanExtended{},
		[][]string{
			{"Workspace", "To Create", "To Update", "To Destroy", "To Replace"},
			{"workspace1", "2", "1", "0", "0"},
			{"workspace2", "0", "0", "1", "1"},
		},
		map[string]map[string][]string{},
		map[string][]string{},
	}

	expected := "| Workspace | To Create | To Update | To Destroy | To Replace |\n"
	expected += "|-|:-:|:-:|:-:|:-:|\n"
	expected += "| workspace1 | 2 | 1 | 0 | 0 |\n"
	expected += "| workspace2 | 0 | 0 | 1 | 1 |\n"
	expected += "\n\n"

	assert.Equal(t, expected, planAnalyzer.GenerateComparisonTable())
}

func TestGenerateLastUpdated(t *testing.T) {
	planAnalyzer := &PlanAnalyzer{
		[]PlanExtended{},
		[][]string{{"Workspace", "To Create", "To Update", "To Destroy", "To Replace"}},
		map[string]map[string][]string{},
		map[string][]string{},
	}

	result := planAnalyzer.GenerateLastUpdated()
	assert.Contains(t, result, "Last Updated: `")
	assert.True(t, strings.HasSuffix(result, "`\n\n"))
}

func TestGenerateReport(t *testing.T) {
	changeSet := map[string][]string{
		Create:  {"resource1"},
		Destroy: {"resource2"},
	}

	planAnalyzer := &PlanAnalyzer{
		[]PlanExtended{},
		[][]string{
			{"Workspace", "To Create", "To Update", "To Destroy", "To Replace"},
			{"workspace1", "1", "0", "1", "0"},
		},
		map[string]map[string][]string{
			"workspace1": changeSet,
		},
		map[string][]string{Create: {"resource1"}},
	}

	result := planAnalyzer.GenerateReport()

	assert.Contains(t, result, "# :clipboard: Terraform Plan Analyzer Report! :clipboard:\n")
	assert.Contains(t, result, "Last Updated: `")
	assert.Contains(t, result, "| Workspace | To Create | To Update | To Destroy | To Replace |\n|-|:-:|:-:|:-:|:-:|\n| workspace1 | 1 | 0 | 1 | 0 |\n\n\n")
	assert.Contains(t, result, "## All Workspaces:pencil2:\n```diff\n+ To Create +\n~ resource1\n```\n\n")
	assert.Contains(t, result, "## Individual Workspaces\n### workspace1 :pencil2::wastebasket:\n```diff\n- To Destroy -\n~ resource2\n```\n\n")
}

func TestGenerateReportEmpty(t *testing.T) {
	planAnalyzer := NewPlanAnalyzer([]PlanExtended{})
	planAnalyzer.ProcessPlans()

	result := planAnalyzer.GenerateReport()

	assert.Contains(t, result, "# :clipboard: Terraform Plan Analyzer Report! :clipboard:\n")
	assert.Contains(t, result, "Last Updated: `")
	assert.Contains(t, result, "| Workspace | :pencil2:To Create | :fountain_pen:To Update | :wastebasket:To Destroy | :scissors:To Replace |\n|-|:-:|:-:|:-:|:-:|\n\n\n")
	assert.NotContains(t, result, "## All Workspaces")
	assert.NotContains(t, result, "## Individual Workspaces")
}

func TestProcessPlansComparisonTable(t *testing.T) {
	plans := []PlanExtended{
		{
			Plan:      tfjson.Plan{},
			ChangeSet: []*tfjson.ResourceChange{},
			ToUpdate:  []string{"resource-update"},
			ToCreate:  []string{"resource-create", "resource-create-2"},
			ToDestroy: []string{},
			ToReplace: []string{"resource-replace"},
			Workspace: "workspaceOne",
		},
		{
			Plan:      tfjson.Plan{},
			ChangeSet: []*tfjson.ResourceChange{},
			ToUpdate:  []string{},
			ToCreate:  []string{},
			ToDestroy: []string{"resource-destroy"},
			ToReplace: []string{},
			Workspace: "workspaceTwo",
		},
	}

	planAnalyzer := NewPlanAnalyzer(plans)
	planAnalyzer.ProcessPlans()

	expected := [][]string{
		{"Workspace", ":pencil2:To Create", ":fountain_pen:To Update", ":wastebasket:To Destroy", ":scissors:To Replace"},
		{"workspaceOne", "2", "1", "0", "1"},
		{"workspaceTwo", "0", "0", "1", "0"},
	}
	assert.Equal(t, expected, planAnalyzer.ComparisonTable)
	assert.Equal(t, map[string][]string{
		Create:  {"resource-create", "resource-create-2"},
		Destroy: {},
		Replace: {"resource-replace"},
		Update:  {"resource-update"},
	}, planAnalyzer.Changes["workspaceOne"])
	assert.Equal(t, map[string][]string{
		Create:  {},
		Destroy: {"resource-destroy"},
		Replace: {},
		Update:  {},
	}, planAnalyzer.Changes["workspaceTwo"])
}
