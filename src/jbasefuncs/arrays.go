// -----------------
// Basic, useful functions as existent in other languages are reimplemented here
// -----------------

package jbasefuncs

import ()

// -----------------
// Check if two slices of strings equal each other
// -----------------

func TestEqSliceStrings(a, b []string) bool {

	if a == nil && b == nil {
		return true
	}

	if a == nil || b == nil {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

// -----------------
// Reimplement array_unique, but based on type of slice
// - Thanks: https://www.dotnetperls.com/duplicates-go
// -----------------

func ArrayIntUnique(elements []int) []int {

	// Use map to record duplicates as we find them.
	encountered := map[int]bool{}
	result := []int{}

	for v := range elements {
		if encountered[elements[v]] == false {
			// Record this element as an encountered element.
			encountered[elements[v]] = true
			// Append to result slice.
			result = append(result, elements[v])
		}
	}
	// Return the new slice.
	return result
}

// The same for arrays of strings
func ArrayStringUnique(elements []string) []string {

	// Use map to record duplicates as we find them.
	encountered := map[string]bool{}
	result := []string{}

	for v := range elements {
		if encountered[elements[v]] == false {
			// Record this element as an encountered element.
			encountered[elements[v]] = true
			// Append to result slice.
			result = append(result, elements[v])
		}
	}
	// Return the new slice.
	return result
}

// Reimplementation of PHP's in_array
func InArrayStr(needle string, haystack []string) bool {
	output := false
	for _, value := range haystack {
		if needle == value {
			output = true
		}
	}
	return output
}

// -----------------
// Sum of all number values in an array []int
// -----------------

func ArraySumInt(arr []int) int {
	output := 0
	for _, p := range arr {
		output += p
	}
	return output
}

// -----------------
// Sum of all number values in an array []float64
// -----------------

func ArraySumFloat(arr []float64) float64 {
	output := float64(0)
	for _, p := range arr {
		output += p
	}
	return output
}

func Max(input []int) int {
	var max int
	for _, value := range input {
		if value > max {
			max = value
		}
	}
	return max
}

func MapMaxCount(input map[int][]string) int {
	highest := 0
	for key, value := range input {
		if len(value) > len(input[highest]) {
			highest = key
		}
	}
	return highest
}

// -----------------
// Function to join a slice of strings to a single string
// -----------------

func JoinSlice(joinwith string, list []string) string {
	output := ""
	for _, p := range list {
		output += p + joinwith
	}
	return output[:len(output)-len(joinwith)]
}
