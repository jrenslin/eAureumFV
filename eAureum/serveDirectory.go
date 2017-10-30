// Display folder contents in a table.
package eAureumFV

import (
	"fmt"
	jbasefuncs "github.com/jrenslin/jbasefuncs"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func serveDirectory(w http.ResponseWriter, r *http.Request) {

	folderLocation := r.URL.Query().Get("p") // Parse GET parameter

	// -----------
	// Replace the beginning of the filepath passed via GET with the corresponding folder
	// -----------

	folderNr := strings.Split(folderLocation, "/")[0]
	folderNrInt, err := strconv.Atoi(folderNr)

	// Check for errors in type conversion / invalid folder numbers
	if err != nil {
		fmt.Fprintf(w, "Folder could not be parsed")
		return
	} else if folderNrInt > len(Settings.Folders)-1 {
		fmt.Fprintf(w, "Invalid folder")
		return
	}
	currentBaseDir := Settings.Folders[folderNrInt]
	folderLocation = strings.Replace(folderLocation, folderNr, currentBaseDir, 1)

	// -----------
	// Get folder contents
	// -----------

	var folderContents map[string][]string // Initialize folderContents
	switch {                               // Check for existence of folder before scanning it
	case jbasefuncs.FileExists(folderLocation):
		folderContents = jbasefuncs.ScandirFilesFolders(folderLocation) // Get Folder contents split into files and folders
	default:
		fmt.Fprintf(w, "Invalid folder.")
		return
	}

	// -----------
	// Filling the output variable
	// -----------

	content := "<main>\n"
	content += "<h1>" + filepath.Base(folderLocation) + "</h1>\n"
	content += getTrailHTML(folderLocation, currentBaseDir, folderNr)

	// -----------
	// Print table of files and folders
	// -----------

	content += "\n\n<table>\n"
	content += "<tr><th>Name</th><th>Size</th><th>Last edit</th></tr>\n"
	counter := 1
	for _, file := range folderContents["folders"] { // Loop over folders
		fi, err := os.Stat(file)
		if err != nil { // Print error message and exit function
			fmt.Println("File information on " + filepath.Base(file) + " not accessible.")
			continue
		}

		file = strings.Replace(file, currentBaseDir, "", 1)
		content += "<tr>\n"
		content += "<td class='directory'><a href='./dir?p=" + folderNr + file + "' id='link" + fmt.Sprint(counter) + "'>" + filepath.Base(file) + "</a></td>\n"
		content += "<td></td>\n"
		content += "<td>" + fmt.Sprint(fi.ModTime().Format("2006-01-02 15:04")) + "</td>\n"
		content += "</tr>\n"
		counter++
	}
	for _, file := range folderContents["files"] { // Loop over files
		fi, err := os.Stat(file)
		if err != nil { // Print error message and exit function
			fmt.Println("File information on " + filepath.Base(file) + " not accessible.")
			continue
		}
		fileSize := fi.Size()

		file = strings.Replace(file, currentBaseDir, "", 1)
		content += "<tr>\n"
		content += "<td class='" + jbasefuncs.GetKindOfFile(file) + "'>"
		content += "<a href='./file?p=" + folderNr + file + "' id='link" + fmt.Sprint(counter) + "'>" + filepath.Base(file) + "</a></td>\n"
		content += "<td>" + jbasefuncs.HumanFilesize(fileSize) + "</td>\n"
		content += "<td>" + fmt.Sprint(fi.ModTime().Format("2006-01-02 15:04")) + "</td>\n"
		content += "</tr>\n"
		counter++
	}

	content += "</table>\n"

	content += "</main>"

	// -----------
	// Output
	// -----------

	setHeaders(w, r)

	fmt.Printf(localOutputFormat, time.Now().Format(timeFormat), "Serving table: ", folderLocation)
	htmlPrintPage(w, r, content, "directoryTable", htmlGetMetatags("Directory: "+filepath.Base(folderLocation), "icon", "description", "keywords"))

}
