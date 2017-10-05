// A small file server
//
// Features:
// - Serves from different user-specified directories
// - Offers an easy-to-use navigation for the folders
// -- Completely keyboard-driven usage is possible
// - Offers the files to the browser using the appropriate HTLM5 tags (if there is one)
// - Offers a CBZ viewer
package main

import (
	"./jbasefuncs"
	"./jhtml"
	"./jsonfuncs"
	"archive/zip"
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// -----------------------
// Constants
// -----------------------

const baseLocation = "../"                     // All settings files are located within this directory
const defaultPort = "9090"                     // The default port to serve on (can be overwritten in settings)
const timeFormat = "[2006-01-02 15:04:05] "    // Time format for output in terminal
const localOutputFormat = "%20s %-20s %20s \n" // Determines output in terminal on the computer running this

// -----------------------
// Global variables
// -----------------------

var Settings jsonfuncs.Settings
var FileIndex []string
var FilesByType = make(map[string][]string)

// -----------------------
// Function to check if the list of folders to be served is empty.
// If no, this means that the setup should/can be run.
// -----------------------

func checkForSettings() bool {
	if len(Settings.Folders) == 0 {
		return true
	} else {
		return false
	}

}

// -----------------------
// Ensures that all necessary files and directories are existent. In not, creates them.
// -----------------------

func ensure_working_environment(folder string) {
	jbasefuncs.EnsureDir(folder + "json")
	jbasefuncs.EnsureDir(folder + "css")
	jbasefuncs.EnsureDir(folder + "js")
	jbasefuncs.EnsureDir(folder + "htm")
	jbasefuncs.EnsureJsonList(folder + "json/navigation.json")
	// Create a settings file if none exists yet.
	// Do not set any folders to serve from. Without any set folders, the setup is triggered.
	if jbasefuncs.FileExists(folder+"json/settings.json") == false {
		jbasefuncs.File_put_contents(folder+"json/settings.json", jsonfuncs.ToJson(jsonfuncs.Settings{Port: defaultPort}))
	}
}

// -----------------------
// Serve html pages embedded in the common
// -----------------------

func ServeStaticText(w http.ResponseWriter, r *http.Request) {
	path := strings.Trim(r.URL.Path[1:], "/")

	if jbasefuncs.FileExists("../htm/"+path+".htm") == false {
		fmt.Fprintf(w, "../htm/"+path+".htm")
		return
	}

	content := `

        <main>` + jbasefuncs.File_get_contents("../htm/"+path+".htm") + "</main>"
	jhtml.Print_page(w, r, content, path, jhtml.Get_metatags(path, "icon", "description", "keywords"))
}

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

// -----------------------------------------
// Functions for serving the different pages
// -----------------------------------------

// -----------------------
// Prints the welcome and setup page
// -----------------------

func serveSetup(w http.ResponseWriter, r *http.Request) {

	// The setup page is almost static, hence it can be stored externally easily.
	// Here it's fetched and the port is inserted.
	content := strings.Replace(jbasefuncs.File_get_contents(baseLocation+"htm/setup.htm"), "%s", Settings.Port, 1)

	fmt.Printf(localOutputFormat, time.Now().Format(timeFormat), "Starting setup", "")
	jhtml.Print_page(w, r, content, "setup", jhtml.Get_metatags("Welcome / Setup", "icon", "description", "keywords"))

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

	// Split folders by line
	folders := strings.Split(rawFolders, "\n")
	for key, value := range folders {
		folders[key] = strings.Trim(value, "\r")
	}

	// Write updated settings to global variable
	Settings.Port = port
	Settings.Folders = folders

	runIndexing() // Re-build indexes

	// Store newly set settings and redirect to start page
	fmt.Printf(localOutputFormat, time.Now().Format(timeFormat), "Storing settings", "")
	jbasefuncs.File_put_contents(baseLocation+"json/settings.json", jsonfuncs.ToJson(Settings))
	http.Redirect(w, r, "/", 301)

}

// -----------------------
// Function serving the start page
// - Lists all folders from settings
// - Displays rudimentary statistics
// -----------------------

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

	content += "<h1>*Name*</h1>\n"
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
	jhtml.Print_page(w, r, content, "startPage", jhtml.Get_metatags("Start page", "icon", "description", "keywords"))
}

// -----------------------
// Prints a table of all files of a given kind (e.g. images, plaintext files, etc.)
// -----------------------

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
	jhtml.Print_page(w, r, content, "typeTable", jhtml.Get_metatags("Table: "+strings.Title(selectedType), "icon", "description", "keywords"))

}

// -----------------------
// Display folder contents in a table.
// -----------------------

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
	content += jhtml.GetTrailHTML(Settings, folderLocation, currentBaseDir, folderNr)

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

	fmt.Printf(localOutputFormat, time.Now().Format(timeFormat), "Serving table: ", folderLocation)
	jhtml.Print_page(w, r, content, "directoryTable", jhtml.Get_metatags("Directory: "+filepath.Base(folderLocation), "icon", "description", "keywords"))

}

// -----------------------
// Serves single files
// -----------------------

func serveFile(w http.ResponseWriter, r *http.Request) {

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

	content := "<main>\n"
	content += "<h1>" + filepath.Base(folderLocation) + "</h1>\n"
	content += jhtml.GetTrailHTML(Settings, folderLocation, currentBaseDir, folderNr)

	// -----------
	// Offer option to show preview in full
	// -----------

	switch {
	case fullSized != "":
		content += "<div class='preview fullsized'>\n"
	default:
		content += "<div class='preview'>\n"
	}

	// -----------
	// Show preview pased on file type of file
	// -----------

	displayType := jbasefuncs.GetKindOfFile(folderLocation)
	switch {
	case displayType == "audio":
		content += jhtml.HtmlAudio("/static/" + r.URL.Query().Get("p"))
	case displayType == "video":
		content += jhtml.HtmlVideo("/static/" + r.URL.Query().Get("p"))
	case displayType == "image":
		content += jhtml.HtmlImage("/static/" + r.URL.Query().Get("p"))
	case displayType == "pdf":
		content += jhtml.HtmlPdf("/static/" + r.URL.Query().Get("p"))
	case displayType == "webpage":
		content += jhtml.HtmlWebPage("/static/" + r.URL.Query().Get("p"))
	case displayType == "plaintext":
		content += jhtml.HtmlPlaintext("/static/"+r.URL.Query().Get("p"), folderLocation)
	case displayType == "code":
		content += jhtml.HtmlCode("/static/"+r.URL.Query().Get("p"), folderLocation)
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
		content += `
	  	<script>
		if(!window.location.hash) {
			window.location.href = "./file?p=` + r.URL.Query().Get("p") + `&offset=` + fmt.Sprint(offset) + `#page0";
		}
		</script>
		`

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
	// Serve output
	// -----------

	fmt.Printf(localOutputFormat, time.Now().Format(timeFormat), "Serving file: ", folderLocation)
	jhtml.Print_page(w, r, content, "file", jhtml.Get_metatags("File: "+filepath.Base(folderLocation), "icon", "description", "keywords"))

}

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

// -----------------------
// Serves index of all files in json
// -----------------------

func serveFileIndex(w http.ResponseWriter, r *http.Request) {
	var output []string

	w.Header().Set("Content-Type", "application/json")
	for _, file := range FileIndex {
		for i, folder := range Settings.Folders {
			switch {
			case strings.HasPrefix(file, folder):
				output = append(output, strings.Replace(file, folder, fmt.Sprint(i), 1))
				break
			}
		}
	}
	w.Header().Set("Etag", `"index"`)
	w.Header().Set("Cache-Control", "max-age=86400") // 86400 = 1 day
	fmt.Fprintf(w, jsonfuncs.ToJson(output))
}

// -----------------------
// Init:
// - Make sure that all important folders are available.
// - Load settings
// - Build indexes of all files and of all files by their types
// -----------------------

func init() {

	fmt.Printf(localOutputFormat, time.Now().Format(timeFormat), "Initializing ... ", "")

	ensure_working_environment(baseLocation)
	Settings = jsonfuncs.DecodeSettings(baseLocation + "json/settings.json") // Settings
	runIndexing()                                                            // Load indexes

}

// -----------------------
// Main
// -----------------------

func main() {

	fmt.Printf(localOutputFormat, time.Now().Format(timeFormat), "Starting ... ", "")

	// Bind callable URLs
	http.HandleFunc("/", serveStartPage)                  // Serve startpage
	http.HandleFunc("/dir", serveDirectory)               // Serve directory table
	http.HandleFunc("/file", serveFile)                   // Serve page for specific files
	http.HandleFunc("/zip", serveZipContents)             // Serve a file out of a ZIP archive
	http.HandleFunc("/type", serveFileTypeTable)          // Serve table of files based on file types
	http.HandleFunc("/storeSettings", serveStoreSettings) // Serve page for storing settings (ATM restricted to initial setup)
	http.HandleFunc("/cheatSheet/", ServeStaticText)      // Serve cheat sheet as an html page embedded into the common layout
	http.HandleFunc("/index", serveFileIndex)             // Serve simple index of all files in json

	http.HandleFunc("/css/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../"+r.URL.Path[1:])
	})
	http.HandleFunc("/js/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../"+r.URL.Path[1:])
	})

	// Serve folders specified in the settings
	for _, value := range Settings.Folders {

		// Write contents of values taken from the loop to loop-specific variables.
		// Without this, the values would change as the loop progresses and
		//   wrong folders' contents would be bound to a given URL path.
		var key string
		var folder string
		for i, f := range Settings.Folders {
			if f == value {
				key = fmt.Sprint(i)
				folder = fmt.Sprint(f)
			}
		}
		http.HandleFunc("/static/"+key+"/", func(w http.ResponseWriter, r *http.Request) {
			fmt.Println(strings.Replace(r.URL.Path, "/static/"+key, folder, 1))
			http.ServeFile(w, r, strings.Replace(r.URL.Path, "/static/"+key, folder, 1))
		})
	}

	// Set port to listen on
	err := http.ListenAndServe(":"+Settings.Port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
