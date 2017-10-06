// File types
package jbasefuncs

import (
	"path/filepath"
	"strings"
)

// Lists file extensions by their use case
var FileTypes = map[string][]string{
	"audio":      []string{".mp3", ".m4a", ".ogg"},
	"video":      []string{".mp4", ".webm"},
	"image":      []string{".gif", ".jpg", ".jpeg", ".png", ".bmp"},
	"pdf":        []string{".pdf"},
	"webpage":    []string{".htm", ".html"},
	"plaintext":  []string{".txt"},
	"code":       []string{".md", ".org", ".py", ".php", ".tex"},
	"compressed": []string{".zip", ".rar", ".7z", ".7zip", ".cbr"},
	"comic":      []string{".cbz"},
}

// Returns the file type as described above of a given file
func GetKindOfFile(filename string) string {

	extension := strings.ToLower(filepath.Ext(filename))
	for output, extensions := range FileTypes {
		if InArrayStr(extension, extensions) {
			return output
		}
	}
	return "other"
}
