// Serves the start page
// - Lists all folders from settings
// - Displays rudimentary statistics
package eAureumFV

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func serveStartPage(w http.ResponseWriter, r *http.Request) {

	// -----------
	// Check if the setup needs to be run
	// -----------
	switch {
	case checkForSettings():
		serveSetup(w, r)
		return // Stop function execution if the setup runs
	}

	// -----------
	// Start filling output variable (content)
	// -----------

	content := "<main>\n"

	content += "<h1>eAureumFV</h1>\n"
	content += "<p class='trail'><a href='/' id='link0'>/</a></p>\n"

	content += "<ul class='tiles'>\n"
	for folderNr, folder := range Settings.Folders { // Loop over folders and list them up
		content += "<li>\n"
		content += "<a class='directory' id='link" + fmt.Sprint(folderNr+1) + "' href='/dir?p=" + fmt.Sprint(folderNr) + "'>" + filepath.Base(folder) + "</a>\n"
		content += "</li>\n"
	}
	content += "</ul>\n"

	content += "</main>\n"

	// -----------
	// Write rudimentary "statistics"
	// -----------

	content += "<section>\n"
	content += "<h2>Numbers</h2>\n"
	content += "<div class='tiled'>"

	content += `<div>
	<dl>
	  <dt>Number of all files</dt><dd>` + fmt.Sprint(len(FileIndex)) + `</dd>
	</dl>
	</div>
	`

	// Bar chart on distribution of file types

	content += "<div>\n<div class='barChart'>\n"
	total := 0
	for _, files := range FilesByType {
		total += len(files)
	}

	for name, files := range FilesByType {
		content += "<a href='/type?q=" + name + "&offset=0'>"
		content += "<span class='' style='height: " + strconv.FormatFloat(float64(100)/float64(total)*float64(len(files)), 'f', 3, 64) + "%%;'></span>\n"
		content += "<span>" + fmt.Sprint(len(files)) + "</span><span>" + strings.Title(name) + "</span>\n"
		content += "</a>\n"
	}
	content += "</div>\n</div>\n"

	content += "</div>\n"
	content += "</section>\n"

	// -----------
	// Serve output
	// -----------

	fmt.Printf(localOutputFormat, time.Now().Format(timeFormat), "Serving startpage: ", "")
	htmlPrintPage(w, r, content, "startPage", htmlGetMetatags("Start page", "icon", "description", "keywords"))
}
