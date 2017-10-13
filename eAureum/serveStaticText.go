// Serve html pages embedded in the common layout
package eAureumFV

import (
        jbasefuncs "github.com/jrenslin/jbasefuncs"
	"fmt"
	"net/http"
	"strings"
)

func serveStaticText(w http.ResponseWriter, r *http.Request) {
	path := strings.Trim(r.URL.Path[1:], "/")

	if jbasefuncs.FileExists(baseLocation+"htm/"+path+".htm") == false {
		fmt.Fprintf(w, baseLocation+"htm/"+path+".htm")
		return
	}

	content := `

        <main>` + jbasefuncs.FileGetContents(baseLocation+"htm/"+path+".htm") + "</main>"
	htmlPrintPage(w, r, content, path, htmlGetMetatags(path, "icon", "description", "keywords"))
}
