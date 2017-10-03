package main

import (
	"./jbasefuncs"
	"./jhtml"
	"./jsonfuncs"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const baseLocation = "../"
const defaultPort = "9090"
const timeFormat = "[2006-01-02 15:04:05] "

var Settings jsonfuncs.Settings

// Function to check if the list of folders to be served is empty.
// If that is the case, run setup.
func checkForSettings(w http.ResponseWriter, r *http.Request) bool {
	if len(Settings.Folders) == 0 {
		return true
	} else {
		return false
	}

}

func ensure_working_environment(folder string) {
	jbasefuncs.EnsureDir(folder + "json")
	jbasefuncs.EnsureDir(folder + "css")
	jbasefuncs.EnsureDir(folder + "js")
	jbasefuncs.EnsureJsonList(folder + "json/navigation.json")
	if jbasefuncs.FileExists(folder+"json/settings.json") == false {
		jbasefuncs.File_put_contents(folder+"json/settings.json", jsonfuncs.ToJson(jsonfuncs.Settings{Port: defaultPort}))
	}
}

func ServeStaticText(w http.ResponseWriter, r *http.Request) {
	path := strings.Trim(r.URL.Path[1:], "/")
	content := `<nav>` + jsonfuncs.Get_navigation("../data", "../data/navigation.json") + `</nav>

<main>` + jbasefuncs.File_get_contents("../data/"+path+".htm") + "</main>"
	jhtml.Print_page(w, r, content, path, jhtml.Get_metatags("Tickets", "icon", "description", "keywords"))
}

// Prints the welcome and setup page
func serveSetup(w http.ResponseWriter, r *http.Request) {

	content := `
<section class="fullpage" id="page1">
  <h1>Welcome</h1>
  <p>This is just a little page I wrote to learn web programming in Go.</p>
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

<script>
  if(!window.location.hash) {
    window.location.href = "./#page1";
  }
</script>
        `
	fmt.Println(time.Now().Format(timeFormat) + "Setup")
	jhtml.Print_page(w, r, content, "setup", jhtml.Get_metatags("Setup", "icon", "description", "keywords"))

}

// Function to store settings sent via a POST request
// First checks if settings have already been saved. Abort if yes, to prevent use of this function for any usage besides the initial setup.
func serveStoreSettings(w http.ResponseWriter, r *http.Request) {

	if checkForSettings(w, r) == false {
		http.Redirect(w, r, "/", 301)
		return
	}

	// Parsing POST variables
	r.ParseForm()                       // Parses the request body
	port := r.Form.Get("port")          // x will be "" if parameter is not set
	rawFolders := r.Form.Get("folders") //

	// Split folders by line
	folders := strings.Split(rawFolders, "\n")
	for key, value := range folders {
		folders[key] = strings.Trim(value, "\r")
	}

	Settings.Port = port
	Settings.Folders = folders
	fmt.Println(time.Now().Format(timeFormat) + "Storing settings")

	jbasefuncs.File_put_contents(baseLocation+"json/settings.json", jsonfuncs.ToJson(Settings))
	http.Redirect(w, r, "/", 301)

}

// Function serving the start page
// Lists all folders from settings
func serveStartPage(w http.ResponseWriter, r *http.Request) {

	// Check if the setup needs to be run
	switch {
	case checkForSettings(w, r):
		serveSetup(w, r)
		return // Stop function execution if the setup runs
	}
	content := "<main>\n"

	content += "<h1>A File Server</h1>\n"
	content += "<p class='trail'><a href='/' id='link0'>/</a></p>\n"

	content += "<ul class='tiles'>\n"
	for folderNr, folder := range Settings.Folders {
		content += "<li>\n"
		content += "<a class='directory' id='link" + fmt.Sprint(folderNr+1) + "' href='/dir?p=" + fmt.Sprint(folderNr) + "'>" + filepath.Base(folder) + "</a>\n"
		content += "</li>\n"
	}
	content += "</ul>\n"

	content += "</main>\n"
	jhtml.Print_page(w, r, content, "startPage", jhtml.Get_metatags("Setup", "icon", "description", "keywords"))
}

// Function to check if the given folder is a subfolder.
// Displays folder contents.
func serveDirectory(w http.ResponseWriter, r *http.Request) {

	folderLocation := r.URL.Query().Get("p")

	// Replace the beginning of the filepath passed via GET
	folderNr := strings.Split(folderLocation, "/")[0]
	folderNrInt, err := strconv.Atoi(folderNr)
	jbasefuncs.Check(err)
	currentBaseDir := Settings.Folders[folderNrInt]
	folderLocation = strings.Replace(folderLocation, folderNr, currentBaseDir, 1)

	var folderContents map[string][]string // Initialize folderContents)

	// Check for invalid folder locations
	if folderLocation == "" {
		http.Redirect(w, r, "/", 301)
	}
	folderContents = jbasefuncs.ScandirFilesFolders(folderLocation)

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
		jbasefuncs.Check(err)

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
		jbasefuncs.Check(err)
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

	fmt.Println(time.Now().Format(timeFormat) + "Serving table: " + folderLocation)
	jhtml.Print_page(w, r, content, "directoryTable", jhtml.Get_metatags("Directory: "+filepath.Base(folderLocation), "icon", "description", "keywords"))

}

func serveFile(w http.ResponseWriter, r *http.Request) {

	folderLocation := r.URL.Query().Get("p")

	// Replace the beginning of the filepath passed via GET
	folderNr := strings.Split(folderLocation, "/")[0]
	folderNrInt, err := strconv.Atoi(folderNr)
	jbasefuncs.Check(err)
	currentBaseDir := Settings.Folders[folderNrInt]
	folderLocation = strings.Replace(folderLocation, folderNr, currentBaseDir, 1)

	// Check folder contents to later offer the option to navigate to the previous / next file
	var folderContents map[string][]string // Initialize folderContents)
	// Check for invalid folder locations
	if folderLocation == "" || jbasefuncs.FileExists(folderLocation) == false {
		http.Redirect(w, r, "/", 301)
	}
	folderContents = jbasefuncs.ScandirFilesFolders(strings.Replace(folderLocation, filepath.Base(folderLocation), "", 1))

	// Check position of the currently selected file
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
	content += "<div class='preview'>\n"
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

	fmt.Println(time.Now().Format(timeFormat) + "Serving file: " + folderLocation)
	jhtml.Print_page(w, r, content, "file", jhtml.Get_metatags("File: "+filepath.Base(folderLocation), "icon", "description", "keywords"))

}

func main() {

	fmt.Println(time.Now().Format(timeFormat) + "Starting ... ")
	ensure_working_environment(baseLocation)
	Settings = jsonfuncs.DecodeSettings(baseLocation + "json/settings.json")

	http.HandleFunc("/", serveStartPage)                  // Serve startpage on
	http.HandleFunc("/dir", serveDirectory)               // Serve directories on
	http.HandleFunc("/file", serveFile)                   // Serve page for specific files
	http.HandleFunc("/storeSettings", serveStoreSettings) // Serve page for storing settings (atm restricted for initial setup)
	http.HandleFunc("/css/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../"+r.URL.Path[1:])
	})
	http.HandleFunc("/js/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../"+r.URL.Path[1:])
	})

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
			fmt.Println("/static/" + key + "/")
			fmt.Println(strings.Replace(r.URL.Path, "/static/"+key, folder, 1))
			http.ServeFile(w, r, strings.Replace(r.URL.Path, "/static/"+key, folder, 1))
		})
	}
	http.HandleFunc("/about/", ServeStaticText)
	err := http.ListenAndServe(":"+Settings.Port, nil) // set listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
