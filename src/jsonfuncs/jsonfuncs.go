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

// -----------------
// Navigation element
// -----------------

type NavSubElement struct {
	Name string `json:"name"`
	Link string `json:"link"`
}

type NavElement struct {
	Name          string          `json:"name"`
	Link          string          `json:"link"`
	Subnavigation []NavSubElement `json:"sub"`
}

// -----------------
// Printing Navigation
// -- Skipping options for reverting navigation back to JSON for now
// -----------------

func Get_navigation(datafolder string, filename string) string {

	file := jbasefuncs.File_get_contents_bytes(filename)

	var data []NavElement
	err := json.Unmarshal(file, &data)

	jbasefuncs.Check(err)

	output := "<ul>\n"
	for _, p := range data {
		output += "  <li>\n"
		output += "    <a href='" + p.Link + "' id='navigation_" + p.Name + "' />" + p.Name + "</a>\n"

		output += "    <ul>\n"
		for _, subElement := range p.Subnavigation {
			output += "      <li>\n"
			output += "        <a href='" + subElement.Link + "' id='navigation_" + subElement.Name + "' />" + subElement.Name + "</a>\n"
			output += "      </li>\n"
		}
		output += "    </ul>\n"

		output += "  </li>\n"
	}
	output += "</ul>\n"

	return output

}
