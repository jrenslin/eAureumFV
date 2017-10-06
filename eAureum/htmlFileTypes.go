package eAureumFV

import (
	"../jbasefuncs"
	"strings"
)

// -------------------------
// Filetype-specific outputs
//
// Each function takes one or two parameters.
// Parameter 1 (src) : The path of the publicly displayed source. Needs to be accepted by the client.
// Parameter 2 (file): The actual filepath of the file. Only needed if some further interaction with the file is necessary.
// -------------------------

func HtmlAudio(src string) string { // HTML output for displaying an audio file
	fileSplit := strings.Split(src, ".")
	fileType := strings.ToLower(fileSplit[len(fileSplit)-1])

	if fileType == "mp3" { // With the other formats, the file type equals the MIME type. MP3 has MIME type mpeg.
		fileType = "mpeg"
	}

	return `
        <audio controls>
          <source src="` + src + `" type="audio/` + fileType + `">
          Your browser does not support the audio tag.<br /><a href="` + src + `">Open the file.</a>
        </audio>
        `
}

func HtmlVideo(src string) string { // HTML output for displaying a video file
	fileSplit := strings.Split(src, ".")
	fileType := strings.ToLower(fileSplit[len(fileSplit)-1])

	return `
        <video controls>
          <source src="` + src + `" type="video/` + fileType + `">
          Your browser does not support the video tag.<br /><a href="` + src + `">Open the file.</a>
        </video>
        `
}
func HtmlImage(src string) string { // HTML output for displaying an image
	return "<img src='" + src + "' />"
}
func HtmlPdf(src string) string { // HTML output for displaying html files in a frame
	return "<object data='" + src + "' type='application/pdf'><a href='" + src + "'>PDF file.</a></object>"
}
func HtmlWebPage(src string) string { // HTML output for displaying html files in a frame
	return "<iframe src='" + src + "'></iframe>"
}
func HtmlPlaintext(src string, file string) string { // HTML output for displaying plain text files in a frame
	return "<div class='plaintextPreview'>" + jbasefuncs.FileGetContents(file) + "</div>"
}
func HtmlCode(src string, file string) string { // HTML output for displaying code in a frame
	// For css line numbering, the output needs to be split into single lines
	fileContent := jbasefuncs.FileGetContents(file)
	var output string
	for _, line := range strings.Split(fileContent, "\n") {
		line = strings.Trim(line, "\r")
		output += "<span>" + line + "</span>\n"
	}
	return "<code>" + output + "</code>"
}
