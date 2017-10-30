// Serves index of all files in json
package eAureumFV

import (
	"fmt"
	"net/http"
	"strings"
)

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
	setHeaders(w, r)
	fmt.Fprintf(w, ToJson(output))
}
