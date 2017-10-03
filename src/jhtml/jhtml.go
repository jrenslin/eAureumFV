package jhtml

import (
	"../jbasefuncs"
	"../jsonfuncs"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Returns an HTML string containing the
func Get_metatags(title string, icon string, description string, keywords string) string {
	// Read file

	metatags := `
	<title>` + title + `</title>
	<meta http-equiv="content-type" content="text/html; charset=utf-8" />
	<meta name="viewport" content="width=480, initial-scale=0.7" />
	<link rel="stylesheet" type="text/css" href="css/main.css" />
`

	return metatags
}

// Print_page makes the program send the HTML page
// It is primarily a shortcut to reduce redundancy in writing HTML heads and their contents.
func Print_page(w http.ResponseWriter, r *http.Request, content string, id string, metatags string) {
	head := `<!DOCTYPE html>
<html>
<head>
` + metatags + `
</head>
<body id="` + id + `">

<nav>` + jsonfuncs.Get_navigation("../json", "../json/navigation.json") + `</nav>


`

	fmt.Fprintf(w, head)    // Send start of structure and metatags
	fmt.Fprintf(w, content) // Send main content

	fmt.Fprintf(w, "\n<script type='text/javascript' src='js/bindCursors.js'></script>") // Add js to the bottom of the body
	fmt.Fprintf(w, "\n\n</body>\n</html>")                                               // Finish the page
}

// -------------------------
// Functions for printing the path trail.
// -------------------------

// Function to split a filepath into it's folders and return them as an array
func GetTrail(path string, basePath string) []string {
	applicablePath := path[len(basePath):]       // Remove the base path.
	output := strings.Split(applicablePath, "/") // Split remaining path
	return output
}

// Function to form the trail  in HTML
func GetTrailHTML(settings jsonfuncs.Settings, folderLocation string, currentBaseDir string, folderNo string) string {

	// Check if the current folderLocation variable leads to a file or folder
	file, err := os.Stat(folderLocation)
	jbasefuncs.Check(err)
	isDir := file.IsDir() // Write to variable to not need to execute function all the time

	folderNoInt, _ := strconv.Atoi(folderNo)                      // Get Int from the folder no.
	trailelements := GetTrail(folderLocation, currentBaseDir)[1:] // Get elements to loop over
	trailLink := folderNo                                         // Traillinks needs to be initialized outsize of the loop

	// Start filling output var
	// The main directory is always linked in the same way
	output := "<p class='trail'><a href='/' id='link0'>/</a>"
	fmt.Println(len(trailelements))
	if len(trailelements) == 1 { // The keybinding for going up (CTRL+Up) takes precedence over the general one for second level dirs.
		output += "<a href='/dir?p=" + folderNo + "' id='goUp'>" + filepath.Base(settings.Folders[folderNoInt]) + "</a>"
	} else {
		output += "<a href='/dir?p=" + folderNo + "' id='link0ctrl'>" + filepath.Base(settings.Folders[folderNoInt]) + "</a>"
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

// -------------------------
// Filetype-specific outputs
//
// Each function takes one or two parameters.
// Parameter 1 (src) : The path of the publicly displayed source. Needs to be accepted by the client.
// Parameter 2 (file): The actual filepath of the file. Only needed if some further interaction with the file is necessary.
// -------------------------

func HtmlAudio(src string) string { // HTML output for displaying an audio file
	fileSplit := strings.Split(src, ".")
	fileType := strings.ToLower(fileSplit[len(fileSplit)-1])

	if fileType == "mp3" { // With the other formats, the file type equals the MIME type. MP3 has MIME type mpeg.
		fileType = "mpeg"
	}

	return `
        <audio controls>
          <source src="` + src + `" type="audio/` + fileType + `">
          Your browser does not support the audio tag.<br /><a href="` + src + `">Open the file.</a>
        </audio>
        `
}

func HtmlVideo(src string) string { // HTML output for displaying a video file
	fileSplit := strings.Split(src, ".")
	fileType := strings.ToLower(fileSplit[len(fileSplit)-1])

	return `
        <video controls>
          <source src="` + src + `" type="video/` + fileType + `">
          Your browser does not support the video tag.<br /><a href="` + src + `">Open the file.</a>
        </video> 
        `
}
func HtmlImage(src string) string { // HTML output for displaying an image
	return "<img src='" + src + "' />"
}
func HtmlPdf(src string) string { // HTML output for displaying html files in a frame
	return "<object data='" + src + "' type='application/pdf'><a href='" + src + "'>PDF file.</a></object>"
}
func HtmlWebPage(src string) string { // HTML output for displaying html files in a frame
	return "<iframe src='" + src + "'></iframe>"
}
func HtmlPlaintext(src string, file string) string { // HTML output for displaying plain text files in a frame
	return "<div class='plaintextPreview'>" + jbasefuncs.File_get_contents(file) + "</div>"
}
func HtmlCode(src string, file string) string { // HTML output for displaying code in a frame
	// For css line numbering, the output needs to be split into single lines
	fileContent := jbasefuncs.File_get_contents(file)
	var output string
	for _, line := range strings.Split(fileContent, "\n") {
		line = strings.Trim(line, "\r")
		output += "<span>" + line + "</span>\n"
	}
	return "<code>" + output + "</code>"
}
