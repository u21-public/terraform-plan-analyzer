package plananalyzer

import (
	"testing"

	tfjson "github.com/hashicorp/terraform-json"
	"github.com/stretchr/testify/assert"
)

func TestSummarize_empty(t *testing.T) {
	pa := NewPlanAnalyzer([]PlanExtended{})
	pa.ProcessPlans()
	s := pa.Summarize()
	assert.Equal(t, Summary{}, s, "empty plans: all counts and WorkspaceCount should be zero")
}

// TestSummarize_destroyOnly verifies that a single pure-delete plan produces
// Destroyed=1 and all other counts=0 (V1).
func TestSummarize_destroyOnly(t *testing.T) {
	plans := []PlanExtended{
		{
			Plan:      tfjson.Plan{},
			ChangeSet: []*tfjson.ResourceChange{},
			ToUpdate:  []string{},
			ToCreate:  []string{},
			ToDestroy: []string{"module.example.aws_s3_bucket.bucket"},
			ToReplace: []string{},
			Workspace: "ws-destroy",
		},
	}
	pa := NewPlanAnalyzer(plans)
	pa.ProcessPlans()
	s := pa.Summarize()
	assert.Equal(t, 0, s.Created)
	assert.Equal(t, 0, s.Modified)
	assert.Equal(t, 1, s.Destroyed, "V1: pure delete must count toward Destroyed")
	assert.Equal(t, 0, s.Replaced)
	assert.Equal(t, 1, s.WorkspaceCount, "V9: WorkspaceCount must equal number of plans")
}

// TestSummarize_replaceOnly_destroyedZero verifies that a create-before-destroy
// replacement produces Replaced=1 and Destroyed=0, enforcing V2.
func TestSummarize_replaceOnly_destroyedZero(t *testing.T) {
	plans := []PlanExtended{
		{
			Plan:      tfjson.Plan{},
			ChangeSet: []*tfjson.ResourceChange{},
			ToUpdate:  []string{},
			ToCreate:  []string{},
			ToDestroy: []string{},
			ToReplace: []string{"module.example.aws_s3_bucket.bucket"},
			Workspace: "ws-replace",
		},
	}
	pa := NewPlanAnalyzer(plans)
	pa.ProcessPlans()
	s := pa.Summarize()
	assert.Equal(t, 1, s.Replaced)
	assert.Equal(t, 0, s.Destroyed, "V2: replacement must NOT count as a destroy")
	assert.Equal(t, 1, s.WorkspaceCount)
}

// TestSummarize_multiWorkspace verifies that counts are summed across all plans
// and ModuleCount reflects the total number of workspaces (V9).
func TestSummarize_multiWorkspace(t *testing.T) {
	plans := []PlanExtended{
		{
			Plan:      tfjson.Plan{},
			ChangeSet: []*tfjson.ResourceChange{},
			ToUpdate:  []string{"res-update-1"},
			ToCreate:  []string{"res-create-1", "res-create-2"},
			ToDestroy: []string{"res-destroy-1"},
			ToReplace: []string{},
			Workspace: "ws-one",
		},
		{
			Plan:      tfjson.Plan{},
			ChangeSet: []*tfjson.ResourceChange{},
			ToUpdate:  []string{},
			ToCreate:  []string{"res-create-3"},
			ToDestroy: []string{"res-destroy-2", "res-destroy-3"},
			ToReplace: []string{"res-replace-1"},
			Workspace: "ws-two",
		},
	}
	pa := NewPlanAnalyzer(plans)
	pa.ProcessPlans()
	s := pa.Summarize()
	assert.Equal(t, 3, s.Created, "created: 2+1")
	assert.Equal(t, 1, s.Modified, "modified: 1+0")
	assert.Equal(t, 3, s.Destroyed, "destroyed: 1+2")
	assert.Equal(t, 1, s.Replaced, "replaced: 0+1")
	assert.Equal(t, 2, s.WorkspaceCount, "V9: two workspaces")
}
