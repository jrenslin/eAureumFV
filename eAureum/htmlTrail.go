// Printing the path trail.
package eAureumFV

import (
        jbasefuncs "github.com/jrenslin/jbasefuncs"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Function to split a filepath into it's folders and return them as an array
func getTrail(path string, basePath string) []string {
	applicablePath := path[len(basePath):]       // Remove the base path.
	output := strings.Split(applicablePath, "/") // Split remaining path
	return output
}

// Function to form the trail  in HTML
func getTrailHTML(folderLocation string, currentBaseDir string, folderNo string) string {

	// Check if the current folderLocation variable leads to a file or folder
	file, err := os.Stat(folderLocation)
	jbasefuncs.Check(err)
	isDir := file.IsDir() // Write to variable to not need to execute function all the time

	folderNoInt, _ := strconv.Atoi(folderNo)                      // Get Int from the folder no.
	trailelements := getTrail(folderLocation, currentBaseDir)[1:] // Get elements to loop over
	trailLink := folderNo                                         // Traillinks needs to be initialized outsize of the loop

	// Start filling output var
	// The main directory is always linked in the same way
	output := "<p class='trail'><a href='/' id='link0'>/</a>"
	fmt.Println(len(trailelements))
	if len(trailelements) == 1 { // The keybinding for going up (CTRL+Up) takes precedence over the general one for second level dirs.
		output += "<a href='/dir?p=" + folderNo + "' id='goUp'>" + filepath.Base(Settings.Folders[folderNoInt]) + "</a>"
	} else {
		output += "<a href='/dir?p=" + folderNo + "' id='link0ctrl'>" + filepath.Base(Settings.Folders[folderNoInt]) + "</a>"
	}

	for _, f := range trailelements {
		trailLink = trailLink + "/" + f
		switch {
		// If the current value is the directory above the provided one, add id="goUp"
		case isDir && len(trailelements) > 1 && f == trailelements[len(trailelements)-2]:
			output += "<a href='/dir?p=" + trailLink + "' id='goUp'>" + f + "</a>"
		// If the given filepath leads to a file and this is the dir above it, add id="goUp"
		case isDir == false && len(trailelements) > 1 && f == trailelements[len(trailelements)-2]:
			output += "<a href='/dir?p=" + trailLink + "' id='goUp'>" + f + "</a>"
		// If the given filepath leads to a file, link it as a file
		case isDir == false && f == trailelements[len(trailelements)-1]:
			output += "<a href='/file?p=" + trailLink + "'>" + f + "</a>"
		default:
			output += "<a href='/dir?p=" + trailLink + "'>" + f + "</a>"
		}
	}

	output += "</p>"
	return output
}
