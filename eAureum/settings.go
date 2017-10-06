// Function to check if the list of folders to be served is empty.
// If no, this means that the setup should/can be run.
package eAureumFV

import (
	"../jbasefuncs"
	"encoding/json"
)

type SettingsType struct {
	Port    string   `json:"port"`
	Folders []string `json:"folders"`
}

// Function for decoding the folder list.
func decodeSettings(filename string) SettingsType {
	file := jbasefuncs.FileGetContentsBytes(filename)

	var data SettingsType
	err := json.Unmarshal(file, &data)
	jbasefuncs.Check(err)

	return data
}

func checkForSettings() bool {
	if len(Settings.Folders) == 0 {
		return true
	} else {
		return false
	}

}
