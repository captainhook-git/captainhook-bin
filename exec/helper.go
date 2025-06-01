package exec

import (
	"os"
	"strings"
)

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

// isSymlink checks if a file is a symlink
func isSymlink(path string) (bool, error) {
	info, err := os.Lstat(path)
	if err != nil {
		return false, err
	}
	return info.Mode()&os.ModeSymlink != 0, nil
}
