package PlanAnalyzer

func getEmojis(changeSet map[string][]string) string {
	emojis := ""
	for action, resources := range changeSet {
		if len(resources) > 0 {
			emojis = emojis + EmojiMap[action]
		}
	}
	return emojis
}

func getGitDiff(action string) string {
	return GitDiffMap[action]
}
