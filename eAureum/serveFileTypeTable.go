// Prints a table of all files of a given kind (e.g. images, plaintext files, etc.)
package eAureumFV

import (
        jbasefuncs "github.com/jrenslin/jbasefuncs"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func serveFileTypeTable(w http.ResponseWriter, r *http.Request) {

	selectedType := r.URL.Query().Get("q") // Parse GET parameter

	offsetStr := r.URL.Query().Get("offset") // Get offset from GET parameters
	if offsetStr == "" {
		offsetStr = "0"
	}
	offset, err := strconv.Atoi(offsetStr) // Parse offset
	if err != nil {
		fmt.Fprintf(w, "Invalid offset")
		return
	}

	// -----------
	// Write content / output variable
	// -----------

	content := "<main>\n"

	content += "<h1>" + strings.Title(selectedType) + "</h1>\n"
	content += "<p class='trail'><a href='/' id='link0'>/</a></p>"

	// -----------
	// Print table of files and folders
	// -----------

	content += "\n\n<table>\n"
	content += "<tr><th>Name</th><th>Size</th><th>Last edit</th></tr>\n"
	counter := 1
	for i, file := range FilesByType[selectedType] { // Loop over files

		switch { // Only show ten files at a time
		case i < offset || i > offset+10:
			continue
		}

		// Get file information
		fi, err := os.Stat(file)
		if err != nil {
			fmt.Println("File information on " + filepath.Base(file) + " not accessible.")
			continue
		}
		fileSize := fi.Size()

		// Check which of the served directories this file is in and add the corresponding link
		for j, availableFolder := range Settings.Folders {
			file = strings.Replace(file, availableFolder, fmt.Sprint(j), 1)
		}

		content += "<tr>\n"
		content += "<td class='" + jbasefuncs.GetKindOfFile(file) + "'>"
		content += "<a href='./file?p=" + file + "' id='link" + fmt.Sprint(counter) + "'>" + filepath.Base(file) + "</a></td>\n"
		content += "<td>" + jbasefuncs.HumanFilesize(fileSize) + "</td>\n"
		content += "<td>" + fmt.Sprint(fi.ModTime().Format("2006-01-02 15:04")) + "</td>\n"
		content += "</tr>\n"
		counter++
	}

	content += "</table>\n"

	content += "<p class='offsetswitchers'>\n"
	content += "<span>" + fmt.Sprint(offset) + " / " + fmt.Sprint(len(FilesByType[selectedType])) + "</span>\n"

	// -----------
	// Print options to switch to next or previous batch / change offset
	// -----------

	if offset >= 10 {
		content += "<a href='/type?q=" + r.URL.Query().Get("q") + "&offset=" + fmt.Sprint(offset-10) + "' rel='prev' id='prev' >" + fmt.Sprint(offset-10) + "</a>\n"
	}
	if offset+10 < len(FilesByType[selectedType]) {
		content += "<a href='/type?q=" + r.URL.Query().Get("q") + "&offset=" + fmt.Sprint(offset+10) + "' rel='next' id='next' >" + fmt.Sprint(offset+10) + "</a>\n"
	}
	content += "</p>\n"

	content += "</main>"

	// -----------
	// Serve output
	// -----------

	fmt.Printf(localOutputFormat, time.Now().Format(timeFormat), "Serving table based on type: ", selectedType)
	htmlPrintPage(w, r, content, "typeTable", htmlGetMetatags("Table: "+strings.Title(selectedType), "icon", "description", "keywords"))

}
