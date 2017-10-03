// -----------------
// JSON-based backend for goclitr.
// -----------------
package jsonfuncs

import (
	"../jbasefuncs"
	"encoding/json"
)

// ------------------------------------------------
// Set functions for different types of JSON files.
// ------------------------------------------------

func ToJson(p interface{}) string {
	bytes, err := json.Marshal(p)
	jbasefuncs.Check(err)
	return string(bytes)
}

type Settings struct {
	Port    string   `json:"port"`
	Folders []string `json:"folders"`
}

// Function for decoding the folder list.
func DecodeSettings(filename string) Settings {
	file := jbasefuncs.File_get_contents_bytes(filename)

	var data Settings
	err := json.Unmarshal(file, &data)
	jbasefuncs.Check(err)

	return data
}

func StoreSettings(filename string, settings Settings) bool {

	return true
}

// -----------------
// Navigation element
// -----------------

type Navelement struct {
	Name string `json:"name"`
	Link string `json:"link"`
}

// -----------------
// Printing Navigation
// -- Skipping options for reverting navigation back to JSON for now
// -----------------

func Get_navigation(datafolder string, filename string) string {

	file := jbasefuncs.File_get_contents_bytes(filename)

	var data []Navelement
	err := json.Unmarshal(file, &data)

	jbasefuncs.Check(err)

	output := "\n"
	for _, p := range data {
		output += "    <a href='" + p.Link + "' />" + p.Name + "</a>\n"
	}

	return output

}
