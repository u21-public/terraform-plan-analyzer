package plananalyzer

import (
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
