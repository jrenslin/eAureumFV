// Print error messages
package jbasefuncs

import (
	"fmt"
	"os"
)

func Check(e error) {
	if e != nil {
		fmt.Println(e.Error())
		os.Exit(1)
	}
}

func Die(str string) {
	fmt.Println(str)
	os.Exit(1)
}
