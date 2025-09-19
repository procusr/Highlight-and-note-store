package templates

// getMethodName maps a learning method key to a human readable name.
func getMethodName(key string) string {
	switch key {
	case "applied_laboratory":
		return "Applied Mastery Laboratory"
	case "foundation_discovery":
		return "Foundation Through Discovery"
	case "integration_transfer":
		return "Integration & Transfer System"
	case "mistake_learning_lab":
		return "Mistake-Based Learning Lab"
	case "problem_gymnasium":
		return "Problem-Solving Gymnasium"
	default:
		return "Unknown Method"
	}
}
