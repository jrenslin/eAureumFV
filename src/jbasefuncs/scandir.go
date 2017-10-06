package jbasefuncs

import (
	"os"
	"path/filepath"
	"strings"
)

// -----------------
// Scandir (following PHP's scandir) lists the contents of a folder
// -----------------

// Returns all available files and folders in a directory
func Scandir(folder string) []string {
	// Ensure that the provided filepath *folder* ends with a string
	if strings.HasSuffix(folder, "/") == false {
		folder += "/"
	}
	files, _ := filepath.Glob(folder + "*")
	return files
}

// Returns all available files and folders in a directory, but offers to restrict the search
func ScandirPlus(folder string, selector string) []string {
	// Ensure that the provided filepath *folder* ends with a string
	if strings.HasSuffix(folder, "/") == false {
		folder += "/"
	}
	files, _ := filepath.Glob(folder + selector)
	return files
}

// Returns all available files and folders in a directory. Returning them distinguishing files and folders
func ScandirFilesFolders(folder string) map[string][]string {
	all := Scandir(folder)
	output := map[string][]string{}

	for _, file := range all {
		fileInfo, err := os.Stat(file)
		Check(err)
		if fileInfo.IsDir() {
			output["folders"] = append(output["folders"], file)
		} else {
			output["files"] = append(output["files"], file)
		}
	}
	return output
}
