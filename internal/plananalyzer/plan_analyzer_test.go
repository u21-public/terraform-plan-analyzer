package plananalyzer

import (
	"testing"

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
	expected = expected + "```diff\n+ To Create +\n```\n\n```diff\n- To Destroy -\n~ resource2\n```\n\n"
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
	expected = expected + "```diff\n+ To Create +\n```\n\n```diff\n- To Destroy -\n~ resource2\n```\n\n"
	expected = expected + "### workspace2 :pencil2::wastebasket:\n```diff\n+ To Create +\n```\n\n```diff\n- To Destroy -\n~ resource2\n```\n\n"
	result := planAnalyzer.generateResources()
	assert.Equal(t, expected, result)
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

	expected := "### workspace1 :pencil2::wastebasket:\n```diff\n+ To Create +\n```\n\n```diff\n- To Destroy -\n~ resource2\n```\n\n"
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
