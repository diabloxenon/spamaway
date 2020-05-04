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

// SetToArr removes duplicates from string array
func SetToArr(elements []string) []string {
	seen := map[string]bool{}
	result := []string{}
	for x := range elements {
		if seen[elements[x]] == true {
			// Do not add duplicate
		} else {
			// This element is seen now
			seen[elements[x]] = true
			// Append to result slice
			result = append(result, elements[x])
		}
	}
	return result
}

// Set removes duplicates from string array
func Set(elements []string) WordDict {
	seen := map[string]bool{}
	result := make(WordDict)
	val := 0
	for x := range elements {
		if seen[elements[x]] != true {
			// This element is seen now
			seen[elements[x]] = true
			result[elements[x]] = val
			val++
		}
		// IGNORE: Do not add duplicate
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

// Count an element in array
func Count(s []string, e string) int {
	tmp := 0
	for i := range s {
		if s[i] == e {
			tmp++
		}
	}
	return tmp
}

func MapToArr(m WordDict) []string {
	arr := []string{}
	for k := range m {
		arr = append(arr, k)
	}
	return arr
}

func SpamFamConfusionMatrix(conMat ConMat, label byte, class string) {
	// True Positive, True Negative, False Positive, False Negative
	// FamFam, FamSpam, SpamFam, SpamSpam
	// conMat = [][]int{[]int{0, 0}, []int{0, 0}}
	if label == 49 && class == "spam" {
		// SpamSpam
		conMat[1][1]++
	} else if label == 48 && class == "fam" {
		//FamFam
		conMat[0][0]++
	} else if label == 49 && class == "fam" {
		//SpamFam
		conMat[1][0]++
	} else if label == 48 && class == "spam" {
		//FamSpam
		conMat[0][1]++
	}
}
