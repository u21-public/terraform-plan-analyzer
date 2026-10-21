package plananalyzer

// Summary holds aggregate resource-change counts across all analyzed plans.
// Only pure deletes (Actions.Delete() == ["delete"]) count toward Destroyed;
// replacements (["create","delete"] or ["delete","create"]) and forgets (["forget"])
// are excluded per V1, V2.
type Summary struct {
	Created        int
	Modified       int
	Destroyed      int
	Replaced       int
	WorkspaceCount int
}

// Summarize aggregates counts across all plans in the analyzer.
// WorkspaceCount equals the number of distinct workspaces/plans analyzed (V9).
// Must be called after ProcessPlans() so ToCreate/ToUpdate/ToDestroy/ToReplace
// are populated on each PlanExtended.
func (pa *PlanAnalyzer) Summarize() Summary {
	s := Summary{WorkspaceCount: len(pa.Plans)}
	for _, plan := range pa.Plans {
		s.Created += len(plan.ToCreate)
		s.Modified += len(plan.ToUpdate)
		s.Destroyed += len(plan.ToDestroy)
		s.Replaced += len(plan.ToReplace)
	}
	return s
}
