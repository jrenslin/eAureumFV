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
	"strconv"
	"strings"
	"time"
)

// Declare constants
const baseLocation = "../"                     // All settings files are located within this directory
const defaultPort = "9090"                     // The default port to serve on (can be overwritten in settings)
const timeFormat = "[2006-01-02 15:04:05] "    // Time format for output in terminal
const localOutputFormat = "%20s %-20s %20s \n" // Determines output in terminal on the computer running this

// Initialize settings as a global variable.
var Settings jsonfuncs.Settings

// Function to check if the list of folders to be served is empty.
// If no, this means that the setup should/can be run.
func checkForSettings() bool {
	if len(Settings.Folders) == 0 {
		return true
	} else {
		return false
	}

}

// Ensures that all necessary files and directories are existent. In not, creates them.
func ensure_working_environment(folder string) {
	jbasefuncs.EnsureDir(folder + "json")
	jbasefuncs.EnsureDir(folder + "css")
	jbasefuncs.EnsureDir(folder + "js")
	jbasefuncs.EnsureJsonList(folder + "json/navigation.json")
	// Create a settings file if none exists yet.
	// Do not set any folders to serve from. Without any set folders, the setup is triggered.
	if jbasefuncs.FileExists(folder+"json/settings.json") == false {
		jbasefuncs.File_put_contents(folder+"json/settings.json", jsonfuncs.ToJson(jsonfuncs.Settings{Port: defaultPort}))
	}
}

// Serve html pages embedded in the common
func ServeStaticText(w http.ResponseWriter, r *http.Request) {
	path := strings.Trim(r.URL.Path[1:], "/")
	content := `
        <nav>` + jsonfuncs.Get_navigation("../data", "../data/navigation.json") + `</nav>

        <main>` + jbasefuncs.File_get_contents("../data/"+path+".htm") + "</main>"
	jhtml.Print_page(w, r, content, path, jhtml.Get_metatags("Page", "icon", "description", "keywords"))
}

// Prints the welcome and setup page
func serveSetup(w http.ResponseWriter, r *http.Request) {

	content := `
        <section class="fullpage" id="page1">
          <h1>Welcome to Name</h1>
          <p>*Name* is a small web-based file server. For more information, see the project page on GitHub.</p>
          <p>To start with, some settings are required. Click below to start the setup.</p>
          <a class='buttonlike' href="#page2">Next: Setup</a>
        </section>

        <section class="fullpage" id="page2">
          <h2>Setup</h2>
          <p>First, some settings need to be done.</p>
          <form action="storeSettings" method="POST">
            <label for="port">Port</label>
            <input type="number" name="port" id="port" value="` + Settings.Port + `" />
            <label for="folders">Folders to serve</label>
            <textarea name="folders" id="folders" placeholder="Put one folder path per line."></textarea>
            <button type="submit">Submit</button>
          </form>
        </section>

        <!-- js: If no hash is set in the url, open first page -->
        <script>
          if(!window.location.hash) {
            window.location.href = "./#page1";
          }
        </script>
        `
	fmt.Printf(localOutputFormat, time.Now().Format(timeFormat), "Starting setup", "")
	jhtml.Print_page(w, r, content, "setup", jhtml.Get_metatags("Welcome / Setup", "icon", "description", "keywords"))

}

// Function to store settings sent via a POST request
// First checks if settings have already been saved.
// Abort if yes, to prevent use of this function for any usage besides the initial setup.
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

	// Write settings to global variable
	Settings.Port = port
	Settings.Folders = folders

	fmt.Printf(localOutputFormat, time.Now().Format(timeFormat), "Storing settings", "")

	// Store newly set settings and redirect to start page
	jbasefuncs.File_put_contents(baseLocation+"json/settings.json", jsonfuncs.ToJson(Settings))
	http.Redirect(w, r, "/", 301)

}

// Function serving the start page
// - Lists all folders from settings
func serveStartPage(w http.ResponseWriter, r *http.Request) {

	// Check if the setup needs to be run
	switch {
	case checkForSettings():
		serveSetup(w, r)
		return // Stop function execution if the setup runs
	}

	// Start filling output variable (content)
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
	jhtml.Print_page(w, r, content, "startPage", jhtml.Get_metatags("Start page", "icon", "description", "keywords"))
}

// Function to check if the given folder is a subfolder. Displays folder contents in a table.
func serveDirectory(w http.ResponseWriter, r *http.Request) {

	folderLocation := r.URL.Query().Get("p") // Parse GET parameter

	// Replace the beginning of the filepath passed via GET
	folderNr := strings.Split(folderLocation, "/")[0]
	folderNrInt, err := strconv.Atoi(folderNr)
	if err != nil { // Print error message and exit function
		fmt.Fprintf(w, "Folder could not be parsed")
		return
	}
	currentBaseDir := Settings.Folders[folderNrInt]
	folderLocation = strings.Replace(folderLocation, folderNr, currentBaseDir, 1)

	if folderLocation == "" { // Check for invalid folder locations
		http.Redirect(w, r, "/", 301)
		return
	}
	var folderContents map[string][]string                          // Initialize folderContents
	folderContents = jbasefuncs.ScandirFilesFolders(folderLocation) // Get Folder contents split into files and folders

	// Write content / output
	content := "<main>\n" // Initialize content variable

	// Add folders to content
	content += "<h1>" + filepath.Base(folderLocation) + "</h1>\n"
	content += jhtml.GetTrailHTML(Settings, folderLocation, currentBaseDir, folderNr)

	// Print table of files and folders
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

	fmt.Printf(localOutputFormat, time.Now().Format(timeFormat), "Serving table: ", folderLocation)
	jhtml.Print_page(w, r, content, "directoryTable", jhtml.Get_metatags("Directory: "+filepath.Base(folderLocation), "icon", "description", "keywords"))

}

func serveFile(w http.ResponseWriter, r *http.Request) {

	folderLocation := r.URL.Query().Get("p")

	// Replace the beginning of the filepath passed via GET
	folderNr := strings.Split(folderLocation, "/")[0]
	folderNrInt, err := strconv.Atoi(folderNr)
	if err != nil { // Print error message and exit function
		fmt.Fprintf(w, "Folder could not be parsed")
		return
	}
	currentBaseDir := Settings.Folders[folderNrInt]
	folderLocation = strings.Replace(folderLocation, folderNr, currentBaseDir, 1)

	// Check folder contents to later offer the option to navigate to the previous / next file
	var folderContents map[string][]string // Initialize folderContents)
	// Check for invalid folder locations
	if folderLocation == "" || jbasefuncs.FileExists(folderLocation) == false {
		http.Redirect(w, r, "/", 301)
	}
	folderContents = jbasefuncs.ScandirFilesFolders(strings.Replace(folderLocation, filepath.Base(folderLocation), "", 1))

	// Check position of the currently selected file within the folder
	// Needed later for links to previous and next files.
	var indexInFolderContents int
	for i, f := range folderContents["files"] {
		if f == folderLocation {
			indexInFolderContents = i
		}
	}

	// Start with filling the output varibale (content)
	content := "<main>\n"
	content += "<h1>" + filepath.Base(folderLocation) + "</h1>\n"
	content += jhtml.GetTrailHTML(Settings, folderLocation, currentBaseDir, folderNr)

	// Show preview pased on file type of file
	displayType := jbasefuncs.GetKindOfFile(folderLocation)

	// Offer option to show preview in full
	fullSized := r.URL.Query().Get("fullPreview")
	switch {
	case fullSized != "":
		content += "<div class='preview fullsized'>\n"
	default:
		content += "<div class='preview'>\n"
	}

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

		offsetStr := r.URL.Query().Get("offset") // Get offset from GET parameters
		if offsetStr == "" {
			offsetStr = "0"
		}
		offset, err := strconv.Atoi(offsetStr) // Parse offset
		if err != nil {
			fmt.Fprintf(w, "Invalid offset")
			return
		}

		archiveContents := listZipContents(folderLocation)
		archiveLocation := r.URL.Query().Get("p")

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

		content += "<div class='infoBox'><div>\n" // Begin of info box
		content += "<p><span id='current'>0</span> / <span id='max'>" + fmt.Sprint(len(archiveContents)) + "</span></p>\n"

		content += "<p class='offsetswitchers'>\n"
		// Print options to switch to next or previous batch / change offset
		if offset >= 10 {
			content += "<a class='offsetswitcher' href='/file?p=" + r.URL.Query().Get("p") + "&offset=" + fmt.Sprint(offset-10) + "' id='prevBatch' >" + fmt.Sprint(offset-10) + "</a>\n"
		}
		if offset+10 < len(archiveContents) {
			content += "<a class='offsetswitcher' href='/file?p=" + r.URL.Query().Get("p") + "&offset=" + fmt.Sprint(offset+10) + "' id='nextBatch' >" + fmt.Sprint(offset+10) + "</a>\n"
		}
		content += "</p>\n"
		content += "</div></div>\n"

		// Add javascript to show first image if no hash is set
		content += `
	  	<script>
		if(!window.location.hash) {
			window.location.href = "./file?p=` + r.URL.Query().Get("p") + `&offset=` + fmt.Sprint(offset) + `#page0";
		}
		</script>
		`

	}
	content += "</div>\n"

	// Add navigation buttons to access next/previous file to output variable
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

	fmt.Printf(localOutputFormat, time.Now().Format(timeFormat), "Serving file: ", folderLocation)
	jhtml.Print_page(w, r, content, "file", jhtml.Get_metatags("File: "+filepath.Base(folderLocation), "icon", "description", "keywords"))

}

// Returns a list of all the file names of files within a ZIP
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

// Serve files from within ZIP-compressed files
// TODO Test this with ZIPs containing folders
func serveZipContents(w http.ResponseWriter, r *http.Request) {

	zipLocation := r.URL.Query().Get("p")
	fileNoStr := r.URL.Query().Get("f") // The number of the file within the ZIP file

	fileNo, err := strconv.Atoi(fileNoStr)
	if err != nil {
		fileNo = 0
	}

	// Replace the beginning of the filepath passed via GET
	folderNr := strings.Split(zipLocation, "/")[0]
	folderNrInt, err := strconv.Atoi(folderNr)
	if err != nil { // Print error message and exit function
		fmt.Fprintf(w, "Folder could not be parsed")
		return
	}
	currentBaseDir := Settings.Folders[folderNrInt]
	zipLocation = strings.Replace(zipLocation, folderNr, currentBaseDir, 1)

	switch { // Prevent invalid file numbers
	case fileNo < 0:
		fileNo = 0
	case fileNo > len(listZipContents(zipLocation)):
		fileNo = jbasefuncs.Max([]int{len(listZipContents(zipLocation)) - 1, 0})
	}

	z, err := zip.OpenReader(zipLocation)
	if err != nil { // Print error message and exit function
		fmt.Fprintf(w, "Archive could not be opened.")
		return
	}
	defer z.Close()

	f := z.File[fileNo]

	rc, err := f.Open()
	if err != nil { // Print error message and exit function
		fmt.Fprintf(w, "File in archive could not be parsed")
		return
	}
	buffer := new(bytes.Buffer)
	buffer.ReadFrom(rc)

	w.Write(buffer.Bytes()) // Serve the file
	rc.Close()              // Close the file

}

func main() {

	fmt.Printf(localOutputFormat, time.Now().Format(timeFormat), "Starting ... ", "")
	ensure_working_environment(baseLocation)
	Settings = jsonfuncs.DecodeSettings(baseLocation + "json/settings.json")

	http.HandleFunc("/", serveStartPage)                  // Serve startpage
	http.HandleFunc("/dir", serveDirectory)               // Serve directory table
	http.HandleFunc("/file", serveFile)                   // Serve page for specific files
	http.HandleFunc("/zip", serveZipContents)             // Serve a file out of a ZIP archive
	http.HandleFunc("/storeSettings", serveStoreSettings) // Serve page for storing settings (ATM restricted to initial setup)
	http.HandleFunc("/css/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../"+r.URL.Path[1:])
	})
	http.HandleFunc("/js/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../"+r.URL.Path[1:])
	})

	// Serve folders specified in the settings
	for _, value := range Settings.Folders {
		var key string
		var folder string
		for i, f := range Settings.Folders {
			if f == value {
				key = fmt.Sprint(i)
				folder = fmt.Sprint(f)
			}
		}
		http.HandleFunc("/static/"+key+"/", func(w http.ResponseWriter, r *http.Request) {
			fmt.Println()
			fmt.Println(strings.Replace(r.URL.Path, "/static/"+key, folder, 1))
			http.ServeFile(w, r, strings.Replace(r.URL.Path, "/static/"+key, folder, 1))
		})
	}
	http.HandleFunc("/about/", ServeStaticText)
	err := http.ListenAndServe(":"+Settings.Port, nil) // Set listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
