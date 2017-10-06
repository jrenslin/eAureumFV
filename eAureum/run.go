// Run the server, listen to given port and determine what to serve
package eAureumFV

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

func Run() {

	fmt.Printf(localOutputFormat, time.Now().Format(timeFormat), "Starting ... ", "")

	// Bind callable URLs
	http.HandleFunc("/", serveStartPage)                  // Serve startpage
	http.HandleFunc("/dir", serveDirectory)               // Serve directory table
	http.HandleFunc("/file", serveFilePage)               // Serve page for specific files
	http.HandleFunc("/zip", serveZipContents)             // Serve a file out of a ZIP archive
	http.HandleFunc("/type", serveFileTypeTable)          // Serve table of files based on file types
	http.HandleFunc("/storeSettings", serveStoreSettings) // Serve page for storing settings (ATM restricted to initial setup)
	http.HandleFunc("/cheatSheet/", serveStaticText)      // Serve cheat sheet as an html page embedded into the common layout
	http.HandleFunc("/index", serveFileIndex)             // Serve simple index of all files in json

	http.HandleFunc("/css/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, baseLocation+r.URL.Path[1:])
	})
	http.HandleFunc("/js/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, baseLocation+r.URL.Path[1:])
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
