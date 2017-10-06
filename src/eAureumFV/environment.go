// Ensures that all necessary files and directories are existent. In not, creates them.
package eAureumFV

import (
	"../jbasefuncs"
)

func ensure_working_environment(folder string) {
	jbasefuncs.EnsureDir(folder + "json")
	jbasefuncs.EnsureDir(folder + "css")
	jbasefuncs.EnsureDir(folder + "js")
	jbasefuncs.EnsureDir(folder + "htm")
	jbasefuncs.EnsureJsonList(folder + "json/navigation.json")
	// Create a settings file if none exists yet.
	// Do not set any folders to serve from. Without any set folders, the setup is triggered.
	if jbasefuncs.FileExists(folder+"json/settings.json") == false {
		jbasefuncs.FilePutContents(folder+"json/settings.json", ToJson(SettingsType{Port: defaultPort}))
	}
}
