package utils

import (
	// "fmt"
)

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