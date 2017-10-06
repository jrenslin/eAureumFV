package jbasefuncs

import (
	"fmt"
)

// -----------------
// Struct for describing categories of timespans, e.g. minutes and their relation to seconds.
// Used to create more understandable output than seconds. See function ReadableTime.
// -----------------

type timecorrespondence struct {
	duration   int64  // Duration in seconds
	descriptor string // Abbreviated form (e.g. h for hour)
}

// -----------------
// Function to convert a time difference (e.g. age) to a human-readable time.
// -----------------

func ReadableTime(seconds int64, roundto bool) string { // Parameter roundto is not used, but included for futher use later.
	correspondence := []timecorrespondence{
		timecorrespondence{duration: 78840000, descriptor: "y"}, // Counts the year as 365 days, ignoring leap years.
		timecorrespondence{duration: 1512000, descriptor: "w"},
		timecorrespondence{duration: 216000, descriptor: "d"},
		timecorrespondence{duration: 3600, descriptor: "h"},
		timecorrespondence{duration: 60, descriptor: "m"},
	}
	for _, timespan := range correspondence {
		if seconds > timespan.duration {
			return fmt.Sprint(seconds/timespan.duration) + timespan.descriptor
		}
	}
	return fmt.Sprint(seconds) + "s"
}

// -----------------
// Functions to present prettier output
// -----------------

// Shortens a filesize to a more readable form.
// Rounds to 2 digits.
func HumanFilesize(filesize int64) string {

	fileSize := float32(filesize)
	switch {
	case fileSize > 1099511627776: // 1099511627776 = 1024*1024*1024*1024
		return fmt.Sprintf("%.2f", fileSize/1099511627776) + " TB"
	case fileSize > 1073741824: // 1073741824 = 1024*1024*1024
		return fmt.Sprintf("%.2f", fileSize/1073741824) + " GB"
	case fileSize > 1048576: // 1048576 = 1024 * 1024
		return fmt.Sprintf("%.2f", fileSize/1048576) + " MB"
	case fileSize > 1024:
		return fmt.Sprintf("%.2f", fileSize/1024) + " KB"
	}
	return fmt.Sprint(filesize) + " B"

}
