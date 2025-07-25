package exec

import (
	"os"
	"strings"
)

// isILogicCondition checks if the condition is an "AND" or an "OR" condition
func isLogicCondition(exec string) bool {
	return strings.HasPrefix(strings.ToLower(exec), "captainhook::logic") ||
		strings.HasPrefix(strings.ToLower(exec), "captainhook.logic")
}

// isInternalFunctionality answers if an action should trigger internal CaptainHook functionality
func isInternalFunctionality(exec string) bool {
	return strings.HasPrefix(strings.ToLower(exec), "captainhook::") ||
		strings.HasPrefix(strings.ToLower(exec), "captainhook.")
}

// splitInternalPath is determining the internal functionality to call
// Internal paths consist of two blocks separated by .
//
// Examples:
// - CaptainHook::SOME.FUNCTIONALITY
// - CaptainHook.Branch.EnsureNaming
func splitInternalPath(exec string) []string {
	var prefix string
	if strings.HasPrefix(strings.ToLower(exec), "captainhook::") {
		prefix = "captainhook::"
	} else {
		prefix = "captainhook."
	}
	pathInfo := strings.Replace(strings.ToLower(exec), prefix, "", 1)
	return strings.Split(pathInfo, ".")
}

// isSymlink checks if a file is a symlink
func isSymlink(path string) (bool, error) {
	info, err := os.Lstat(path)
	if err != nil {
		return false, err
	}
	return info.Mode()&os.ModeSymlink != 0, nil
}
