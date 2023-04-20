package PlanAnalyzer

import (
	"sort"
)

func getEmojis(changeSet map[string][]string) string {
	emojis := ""

	for _, action := range SupportedAction {
		resources := changeSet[action]
		if len(resources) > 0 {
			emojis = emojis + EmojiMap[action]
		}
	}
	return emojis
}

func getGitDiff(action string) (string, bool) {
	result, exists := gitDiffMap[action]
	return result, exists
}

func GetSortedWorkspaces(workspaceChangeSet map[string]map[string][]string) []string {
	sortedWorkspaces := make([]string, len(workspaceChangeSet))

	i := 0
	for workspace := range workspaceChangeSet {
		sortedWorkspaces[i] = workspace
		i++
	}
	sort.Slice(sortedWorkspaces[:], func(i, j int) bool {
		return sortedWorkspaces[i] < sortedWorkspaces[j]
	})

	return sortedWorkspaces
}
