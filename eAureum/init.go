// -----------------------
// Init:
// - Make sure that all important folders are available.
// - Load settings
// - Build indexes of all files and of all files by their types
// -----------------------
package eAureumFV

import (
	"fmt"
	"time"
)

func init() {

	fmt.Printf(localOutputFormat, time.Now().Format(timeFormat), "Initializing ... ", "")

	ensure_working_environment(baseLocation)
	Settings = decodeSettings(baseLocation + "json/settings.json") // Settings
	runIndexing()                                                  // Load indexes

}
