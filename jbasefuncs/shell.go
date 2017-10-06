package jbasefuncs

import ()

// -----------------
// Shortcut for handling command line inputs
// -----------------

func HandleCmdInput(args []string, condition []string) bool {
	if len(args) < len(condition) {
		return false
	} else if TestEqSliceStrings(args[:len(condition)], condition) {
		return true
	} else {
		return false
	}
}
