package PlanAnalyzer

const (
	Create  = "Create"
	Destroy = "Destroy"
	Update  = "Update"
	Replace = "Replace"
)

var SupportedAction = []string{Create, Destroy, Update, Replace}

var EmojiMap = map[string]string{
	Create:  ":pencil2:",
	Destroy: ":wastebasket:",
	Update:  ":fountain_pen:",
	Replace: ":scissors:",
	"title": ":clipboard:",
}

var gitDiffMap = map[string]string{
	Create:  "+",
	Destroy: "-",
	Update:  "!",
	Replace: "-",
}
