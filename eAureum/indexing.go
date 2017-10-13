package eAureumFV

import (
        jbasefuncs "github.com/jrenslin/jbasefuncs"
	"archive/zip"
	"regexp"
)

// -----------------------
// Scan directories recursively
// -----------------------

func scandirRecursive(filepath string) []string {
	var output []string
	folderContents := jbasefuncs.ScandirFilesFolders(filepath)
	for _, folder := range folderContents["folders"] {
		filesInSubdirectory := scandirRecursive(folder)
		output = append(output, filesInSubdirectory...)
	}
	output = append(output, folderContents["files"]...)
	return output
}

// -----------------------
// Creates index of all files by running scandirRecursive over each folder specified in the settings
// -----------------------

func indexAllFiles() []string {

	var output []string
	for _, i := range Settings.Folders {
		output = append(output, scandirRecursive(i)...)
	}
	return output

}

// -----------------------
// Returns all files with any of a given list of file extensions
// -----------------------

func searchIndexForFileExtensions(extensions []string) []string {
	var output []string
	for _, extension := range extensions {
		r := regexp.MustCompile(extension + `$`)
		for _, file := range FileIndex {
			if r.MatchString(file) {
				output = append(output, file)
			}
		}
	}
	return output
}

// -----------------------
// Runs both indexing functions and builds the indexes
// -----------------------

func runIndexing() {
	FileIndex = indexAllFiles()                          // Load index of all files
	for name, extensions := range jbasefuncs.FileTypes { // Index files based on extensions / types
		FilesByType[name] = searchIndexForFileExtensions(extensions)
	}

}

// -----------------------
// Returns a list of all the file names of files within a ZIP
// -----------------------

func listZipContents(file string) []string {

	var output []string

	z, err := zip.OpenReader(file)
	// Return empty string if an error occured when opening the ZIP
	if err != nil {
		return []string{}
	}
	defer z.Close()

	for _, f := range z.File {
		output = append(output, f.Name) // Add every file name to output
	}

	return output
}
