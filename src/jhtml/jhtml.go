package jhtml

import (
	"../jsonfuncs"
	"fmt"
	"net/http"
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
// Filetype-specific outputs
// -------------------------

func HtmlImage(src string) string { // HTML output for displaying an image
	return "<img src='" + src + "' />"
}

func HtmlWebPage(src string) string { // HTML output for displaying html files in a frame
	return "<img src='" + src + "' />"
}
