package eAureumFV

import (
	"archive/zip"
	"bytes"
	"fmt"
	jbasefuncs "github.com/jrenslin/jbasefuncs"
	"net/http"
	"strconv"
	"strings"
)

// -----------------------
// Serve files from within ZIP-compressed files
// TODO Test this with ZIPs containing folders
// -----------------------

func serveZipContents(w http.ResponseWriter, r *http.Request) {

	zipLocation := r.URL.Query().Get("p")
	fileNoStr := r.URL.Query().Get("f") // The number of the file within the ZIP file

	// -----------
	// Get actual location of the archive
	// -----------

	// Replace the beginning of the filepath passed via GET with the corresponding folder
	folderNr := strings.Split(zipLocation, "/")[0]
	folderNrInt, err := strconv.Atoi(folderNr)

	// Check for invalid folder no.
	if err != nil {
		fmt.Fprintf(w, "Folder could not be parsed")
		return
	} else if folderNrInt > len(Settings.Folders)-1 {
		fmt.Fprintf(w, "Invalid folder")
		return
	}
	currentBaseDir := Settings.Folders[folderNrInt]
	zipLocation = strings.Replace(zipLocation, folderNr, currentBaseDir, 1)

	// -----------
	// Get Number of file within the ZIP archive
	// -----------

	fileNo, err := strconv.Atoi(fileNoStr)
	if err != nil {
		fileNo = 0
	}

	switch { // Prevent invalid file numbers
	case fileNo < 0:
		fileNo = 0
	case fileNo > len(listZipContents(zipLocation)):
		fileNo = jbasefuncs.Max([]int{len(listZipContents(zipLocation)) - 1, 0})
	}

	// -----------
	// Set headers before any output
	// -----------

	setHeaders(w, r)

	// -----------
	// Read ZIP file
	// -----------

	z, err := zip.OpenReader(zipLocation)
	if err != nil { // Print error message and exit function
		fmt.Fprintf(w, "Archive could not be opened.")
		return
	}
	defer z.Close()

	// -----------
	// Read file from archive into buffer
	// -----------

	f := z.File[fileNo]

	rc, err := f.Open()
	if err != nil { // Print error message and exit function
		fmt.Fprintf(w, "File in archive could not be parsed")
		return
	}
	buffer := new(bytes.Buffer)
	buffer.ReadFrom(rc)

	// -----------
	// Serve contents of buffer
	// -----------

	w.Write(buffer.Bytes()) // Serve the file
	rc.Close()              // Close the file

}
