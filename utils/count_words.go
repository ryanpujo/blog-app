package utils

import (
	"regexp"
	"strings"
)

func CountWords(s string) uint {
	// Use a regular expression to replace all non-alphanumeric characters with spaces
	reg := regexp.MustCompile(`\W`)
	s = reg.ReplaceAllString(s, " ")

	// Split the string into words using white space as the delimiter
	words := strings.Fields(s)

	// Return the number of words
	return uint(len(words))
}
