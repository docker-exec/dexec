package util

import "regexp"

// AddPrefix takes a string slice and returns a new string slice
// with the supplied prefix inserted before every string in the
// original slice.
func AddPrefix(inSlice []string, prefix string) []string {
	var outSlice []string
	for _, option := range inSlice {
		outSlice = append(outSlice, []string{prefix, option}...)
	}
	return outSlice
}

// JoinStringSlices takes an arbitrary number of string slices
// and concatenates them in the order supplied.
func JoinStringSlices(slices ...[]string) []string {
	var outSlice []string
	for _, slice := range slices {
		outSlice = append(outSlice, slice...)
	}
	return outSlice
}

// ExtractFileExtension extracts the extension from a filename. This is defined
// as the remainder of the string after the last '.'.
func ExtractFileExtension(filename string) string {
	patternPermission := regexp.MustCompile(`.*\.(.*):.*`)
	permissionMatch := patternPermission.FindStringSubmatch(filename)
	if len(permissionMatch) > 0 {
		return permissionMatch[1]
	}
	patternFilename := regexp.MustCompile(`.*\.(.*)`)
	return patternFilename.FindStringSubmatch(filename)[1]
}
