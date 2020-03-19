package utils

import (
	"fmt"
	"io"
)

// PrintCmd cmd output
func PrintCmd(w io.Writer, cmd string) {
	fmt.Fprintln(w, cmd)
}

// PrintCmds cmds output
func PrintCmds(w io.Writer, cmds []string) {
	for _, cmd := range cmds {
		fmt.Fprintln(w, cmd)
	}
}

// RemoveDuplicatesString remove duplicate string in slice
func RemoveDuplicatesString(elements []string) []string {
	encountered := map[string]bool{}

	// Create a map of all unique elements.
	for v := range elements {
		encountered[elements[v]] = true
	}

	// Place all keys from the map into a slice.
	result := []string{}
	for key := range encountered {
		result = append(result, key)
	}
	return result
}
