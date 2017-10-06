// Printing an html page into the correct layout
package eAureumFV

import (
	"fmt"
	"net/http"
)

// Returns an HTML string containing the correct contents of <head></head>
func htmlGetMetatags(title string, icon string, description string, keywords string) string {

	metatags := `
	<title>` + title + `</title>
	<meta http-equiv="content-type" content="text/html; charset=utf-8" />
	<meta name="viewport" content="width=480, initial-scale=0.7" />
	<link rel="stylesheet" type="text/css" href="/css/main.css" />
	`

	return metatags
}

// Print_page makes the program send the HTML page
// It is primarily a shortcut to reduce redundancy in writing HTML heads and their contents.
func htmlPrintPage(w http.ResponseWriter, r *http.Request, content string, id string, metatags string) {

	head := `<!DOCTYPE html>
	<html>
	<head>
	` + metatags + `
	</head>
	<body id="` + id + `">

        <div id="search" style="display: none;">
                <div>
                    <h2>Search Files</h2>
                    <input id="searchInput" name="p" />
                    <ul id="searchSelectors"></ul>
                </div>
        </div>

	<nav>` + getNavigation("../json", "../json/navigation.json") + `</nav>
	`

	fmt.Fprintf(w, head)    // Send start of structure and metatags
	fmt.Fprintf(w, content) // Send main content

	fmt.Fprintf(w, "\n<script type='text/javascript' src='/js/complete.js'></script>")    // Add js for completioon in search
	fmt.Fprintf(w, "\n<script type='text/javascript' src='/js/keybindings.js'></script>") // Add js for keybindings
	fmt.Fprintf(w, "\n\n</body>\n</html>")                                                // Finish the page
}
