package exec

import "strings"

// isILogicCondition checks if the condition is an "AND" or an "OR" condition
func isLogicCondition(action string) bool {
	return strings.HasPrefix(strings.ToLower(action), "captainhook::logic")
}

// isInternalFunctionality answers if an action should trigger internal CaptainHook functionality
func isInternalFunctionality(action string) bool {
	return strings.HasPrefix(strings.ToLower(action), "captainhook::")
}

// splitInternalPath is determining the internal functionality to call
// Internal paths consist of two blocks separated by .
//
// Examples:
// - CaptainHook::SOME.FUNCTIONALITY
// - CaptainHook::Branch.EnsureNaming
func splitInternalPath(action string) []string {
	actionPath := strings.Split(action, "::")[1]
	return strings.Split(actionPath, ".")
}
