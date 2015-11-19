package util

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

const sanitisedWindowsPathPattern = "/%s%s"

// ReadStdin takes a prompt message and
func ReadStdin(promptMessage string, completeMessage string) []string {
	stat, _ := os.Stdin.Stat()
	isPipe := (stat.Mode() & os.ModeCharDevice) == 0
	if !isPipe && len(promptMessage) > 0 {
		fmt.Fprintln(os.Stderr, promptMessage)
	}
	lines := []string{}
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if !isPipe && len(completeMessage) > 0 {
		fmt.Fprintln(os.Stderr, completeMessage)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(scanner.Err())
	}
	return lines
}

// SanitisePath takes an absolute path as provided by filepath.Abs() and
// makes it ready to be passed to Docker based on the current OS. So far
// the only OS format that requires transforming is Windows which is provided
// in the form 'C:\some\path' but Docker requires '/c/some/path'.
func SanitisePath(path string, platform string) string {
	sanitised := path
	if platform == "windows" {
		windowsPathPattern := regexp.MustCompile("^([A-Za-z]):(.*)")
		match := windowsPathPattern.FindStringSubmatch(path)

		driveLetter := strings.ToLower(match[1])
		pathRemainder := strings.Replace(match[2], "\\", "/", -1)

		sanitised = fmt.Sprintf(sanitisedWindowsPathPattern, driveLetter, pathRemainder)
	}
	return sanitised
}

// RetrievePath takes an array whose first element may contain an overridden
// path and converts either this, or the default of "." to an absolute path
// using Go's file utilities. This is then passed to SanitisedPath with the
// current OS to get it into a Docker ready format.
func RetrievePath(targetDirs []string) string {
	path := "."
	if len(targetDirs) > 0 {
		path = targetDirs[0]
	}
	absPath, _ := filepath.Abs(path)
	return SanitisePath(absPath, runtime.GOOS)
}

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

// WriteFile writes a file.
func WriteFile(filename string, content []byte) {
	if err := ioutil.WriteFile(filename, content, 0644); err != nil {
		log.Fatalf("Unable to write %s\n%q", filename, err)
	}
}

// DeleteFile deletes a file.
func DeleteFile(filename string) {
	if err := os.Remove(filename); err != nil {
		log.Fatalf("Unable to delete %s\n%q", filename, err)
	}
}
