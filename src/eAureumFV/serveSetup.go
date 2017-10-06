package eAureumFV

import (
	"../jbasefuncs"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// -----------------------
// Prints the welcome and setup page
// -----------------------

func serveSetup(w http.ResponseWriter, r *http.Request) {

	// The setup page is almost static, hence it can be stored externally easily.
	// Here it's fetched and the port is inserted.
	content := strings.Replace(jbasefuncs.FileGetContents(baseLocation+"htm/setup.htm"), "%s", Settings.Port, 1)

	fmt.Printf(localOutputFormat, time.Now().Format(timeFormat), "Starting setup", "")
	htmlPrintPage(w, r, content, "setup", htmlGetMetatags("Welcome / Setup", "icon", "description", "keywords"))

}

// -----------------------
// Function to store settings sent via a POST request
// First checks if settings have already been saved.
// Abort if yes, to prevent use of this function for any usage besides the initial setup.
// -----------------------

func serveStoreSettings(w http.ResponseWriter, r *http.Request) {

	if checkForSettings() == false { // Check for the existence of sufficient settings.
		http.Redirect(w, r, "/", 301)
		return
	}

	// Parsing POST variables
	r.ParseForm() // Parse POST variables. Necessary to fetch their contents next.
	port := r.Form.Get("port")
	rawFolders := r.Form.Get("folders")

	rawFolders = strings.Trim(rawFolders, " \n") // Remove trailing lines

	// Split folders by line
	folders := strings.Split(rawFolders, "\n")
	for key, value := range folders {
		folders[key] = strings.Trim(value, "\r")
	}

	// Write updated settings to global variable
	Settings.Port = port
	Settings.Folders = folders

	fmt.Println(baseLocation + "json/settings.json")

	// Store newly set settings and redirect to start page
	fmt.Printf(localOutputFormat, time.Now().Format(timeFormat), "Storing settings", "")
	jbasefuncs.FilePutContents(baseLocation+"json/settings.json", ToJson(Settings))

	runIndexing() // Re-build indexes
	http.Redirect(w, r, "/", 301)

}
