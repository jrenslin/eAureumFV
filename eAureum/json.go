// General json-related functions
package eAureumFV

import (
        jbasefuncs "github.com/jrenslin/jbasefuncs"
	"encoding/json"
)

func ToJson(p interface{}) string {
	bytes, err := json.MarshalIndent(p, "", "    ")
	jbasefuncs.Check(err)
	return string(bytes)
}
