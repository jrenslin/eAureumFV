package jbasefuncs

import (
	"io/ioutil"
	"os"
)

// -----------------
// Function to check if there is a
// -----------------

func FileExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

// -----------------
// The following functions EnsureDir and EnsureJson are used to create an empty JSON file / an empty directy if there is none existent yet at the specified filepath.
// -----------------

func EnsureDir(folder string) {
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		os.Mkdir(folder, 0755)
	}
}

func EnsureJson(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		FilePutContents(path, "{}")
	}
}

func EnsureJsonList(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		FilePutContents(path, "[]")
	}
}

// -----------------
// The following functions are basic ports of PHP's file_get_contents and file_put_contents for local usage. As e.g. parsing JSON requires a byte map, not string, a new function for reading the file contents to that without parsing it to a string has been added.
// -----------------

func FileGetContentsBytes(path string) []byte {
	file, e := ioutil.ReadFile(path)
	Check(e)
	return file
}

func FileGetContents(path string) string {
	file, e := ioutil.ReadFile(path)
	Check(e)
	return string(file)
}

func FilePutContents(path string, contents string) {
	d1 := []byte(contents)
	err := ioutil.WriteFile(path, d1, 0644)
	Check(err)
}
