// Handles the page navigation
// The page navigation is stored in a json file
package eAureumFV

import (
	"../jbasefuncs"
	"encoding/json"
)

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

func getNavigation(datafolder string, filename string) string {

	file := jbasefuncs.FileGetContentsBytes(filename)

	var data []NavElement
	err := json.Unmarshal(file, &data)

	jbasefuncs.Check(err)

	jbasefuncs.FilePutContents("../json/navigation.json", ToJson(data))

	output := "<ul>\n"
	for _, p := range data {
		output += "  <li>\n"
		output += "    <a "
		if p.Link != "" { // Only insert href attribute if there is an aim for it.
			output += "href='" + p.Link + "'"
		}
		output += " id='navigation_" + p.Name + "' />" + p.Name + "</a>\n"

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
