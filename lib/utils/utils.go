package utils

import (
	// "fmt"
	"unicode"
)

// Check the errors
func Check(e error) {
    if e != nil {
        panic(e)
    }
}

// Set removes duplicates from string array
func Set(elements []string) []string{
	seen := map[string]bool{}
	result := []string{}
	for x := range elements {
		if seen[elements[x]] == true{
			// Do not add duplicate
		} else{
			// This element is seen now
			seen[elements[x]] = true
			// Append to result slice
			result = append(result, elements[x])
		}
	}
	return result
}

// IsAlpha Checks if the string is alphabet or not
func IsAlpha(s string) bool {
    for _, r := range s {
        if !unicode.IsLetter(r) {
            return false
        }
    }
    return true
}

// Pop elements out of array
func Pop(s []string, i int) []string {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}