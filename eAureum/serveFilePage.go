// Serves a preview page for single files
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

func serveFilePage(w http.ResponseWriter, r *http.Request) {

	folderLocation := r.URL.Query().Get("p")
	fullSized := r.URL.Query().Get("fullPreview")

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
	// Check folder contents
	// to later offer the option to navigate to the previous / next file.
	// -----------

	var folderContents map[string][]string // Initialize folderContents)

	// Check for invalid folder locations
	if folderLocation == "" || jbasefuncs.FileExists(folderLocation) == false {
		fmt.Fprintf(w, "Invalid filepath")
		return
	}
	folderContents = jbasefuncs.ScandirFilesFolders(strings.Replace(folderLocation, filepath.Base(folderLocation), "", 1))

	// -----------
	// Check position of the currently selected file within the folder.
	// Needed later for links to previous and next files.
	// -----------

	var indexInFolderContents int
	for i, f := range folderContents["files"] {
		if f == folderLocation {
			indexInFolderContents = i
			break
		}
	}

	// -----------
	// Start with filling the output varibale (content)
	// -----------

	var content string
	switch {
	case fullSized != "":
		content += "<main class='fullsized'>\n"
	default:
		content += "<main>\n"
	}
	content += "<h1>" + filepath.Base(folderLocation) + "</h1>\n"
	content += getTrailHTML(folderLocation, currentBaseDir, folderNr)

	// -----------
	// Offer option to show preview in full
	// -----------

	content += "<div class='preview'>\n"

	// -----------
	// Show preview pased on file type of file
	// -----------

	displayType := jbasefuncs.GetKindOfFile(folderLocation)
	switch {
	case displayType == "audio":
		content += HtmlAudio("/static/" + r.URL.Query().Get("p"))
	case displayType == "video":
		content += HtmlVideo("/static/" + r.URL.Query().Get("p"))
	case displayType == "image":
		content += HtmlImage("/static/" + r.URL.Query().Get("p"))
	case displayType == "pdf":
		content += HtmlPdf("/static/" + r.URL.Query().Get("p"))
	case displayType == "webpage":
		content += HtmlWebPage("/static/" + r.URL.Query().Get("p"))
	case displayType == "plaintext":
		content += HtmlPlaintext("/static/"+r.URL.Query().Get("p"), folderLocation)
	case displayType == "code":
		content += HtmlCode("/static/"+r.URL.Query().Get("p"), folderLocation)
	case displayType == "comic": // Display for CBZ files is too specialized to move to html package

		// -----------
		// Handle offsets
		// -----------
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
		// Get contents and display them
		// -----------
		archiveContents := listZipContents(folderLocation)
		archiveLocation := r.URL.Query().Get("p") // Get GET variable p again to use it in links

		if len(archiveContents) == 0 { // Stop if the ZIP is empty or invalid.
			return
		}
		for i, _ := range archiveContents {
			switch { // Only show ten files at a time
			case i < offset || i > offset+10:
				continue
			}
			content += "<img src='/zip?p=" + archiveLocation + "&f=" + fmt.Sprint(i) + "' id='page" + fmt.Sprint(i-offset) + "'/>\n"
		}

		// -----------
		// Print info box displaying the number of the currently display file
		// and offering options to change the offset
		// -----------

		content += "<div class='infoBox'><div>\n" // Begin of info box
		content += "<p><span id='current'>0</span> / <span id='max'>" + fmt.Sprint(len(archiveContents)) + "</span></p>\n"

		// Print options to switch to next or previous batch / change offset
		content += "<p class='offsetswitchers'>\n"
		if offset >= 10 {
			content += "<a class='offsetswitcher' href='/file?p=" + r.URL.Query().Get("p") + "&offset=" + fmt.Sprint(offset-10) + "' id='prevBatch' >" + fmt.Sprint(offset-10) + "</a>\n"
		}
		if offset+10 < len(archiveContents) {
			content += "<a class='offsetswitcher' href='/file?p=" + r.URL.Query().Get("p") + "&offset=" + fmt.Sprint(offset+10) + "' id='nextBatch' >" + fmt.Sprint(offset+10) + "</a>\n"
		}
		content += "</p>\n"
		content += "</div></div>\n"

		// -----------
		// Add javascript to show first image if no hash is set
		// -----------
		content += `<script type='text/javascript' src='/js/demandHash.js'></script>`

	}
	content += "</div>\n"

	// -----------
	// Add navigation links to access next/previous file to output variable
	// -----------

	var folderLinkPrev string
	var folderLinkNext string
	if indexInFolderContents >= 1 {
		folderLinkPrev = strings.Replace(folderContents["files"][indexInFolderContents-1], currentBaseDir, folderNr, 1)
		content += "<a href='./file?p=" + folderLinkPrev + "' id='prev' rel='prev'></a>\n"
	}
	if indexInFolderContents < len(folderContents["files"])-1 {
		folderLinkNext = strings.Replace(folderContents["files"][indexInFolderContents+1], currentBaseDir, folderNr, 1)
		content += "<a href='./file?p=" + folderLinkNext + "' id='next' rel='next'></a>\n"
	}

	content += "</main>\n"

	// -----------
	// File information
	// -----------
	fi, err := os.Stat(folderLocation)
	if err != nil { // Print error message and exit function
		fmt.Println("File information on " + filepath.Base(folderLocation) + " not accessible.")
	}
	fileSize := fi.Size()

	content += "<section id='fileInfo'>\n"

	// File options: Download etc.
	content += "<p id='fileOptions'>\n"
	content += "<a title='Download' href='/static/" + r.URL.Query().Get("p") + "' download>&#x2BC6;</a>\n"
	content += "</p>\n"

	content += "<dl>\n"

	content += "<dt>File Size</dt><dd>" + jbasefuncs.HumanFilesize(fileSize) + "</dd>\n"
	content += "<dt>Last Modified</dt><dd>" + fmt.Sprint(fi.ModTime().Format("2006-01-02 15:04")) + "</dd>\n"

	content += "</dl>"
	content += "</section>\n"

	// -----------
	// Serve output
	// -----------

	setHeaders(w, r)

	fmt.Printf(localOutputFormat, time.Now().Format(timeFormat), "Serving file: ", folderLocation)
	htmlPrintPage(w, r, content, "file", htmlGetMetatags("File: "+filepath.Base(folderLocation), "icon", "description", "keywords"))

}
