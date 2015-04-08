package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

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

// DexecImage consists of the file extension, Docker image name and Docker
// image version to use for a given Docker Exec image.
type DexecImage struct {
	extension string
	image     string
	version   string
}

// LookupImageByExtension is a closure storing a dictionary mapping source
// extensions to the names and versions of Docker Exec images.
var LookupImageByExtension = func() func(string) DexecImage {
	innerMap := map[string]DexecImage{
		"c":      {"c", "c", "1.0.0"},
		"clj":    {"clj", "clojure", "1.0.0"},
		"coffee": {"coffee", "coffee", "1.0.0"},
		"cpp":    {"cpp", "cpp", "1.0.0"},
		"cs":     {"cs", "csharp", "1.0.0"},
		"d":      {"d", "d", "1.0.0"},
		"erl":    {"erl", "erlang", "1.0.0"},
		"fs":     {"fs", "fsharp", "1.0.0"},
		"go":     {"go", "go", "1.0.0"},
		"groovy": {"groovy", "groovy", "1.0.0"},
		"hs":     {"hs", "haskell", "1.0.0"},
		"java":   {"java", "java", "1.0.0"},
		"lisp":   {"lisp", "lisp", "1.0.0"},
		"js":     {"js", "node", "1.0.0"},
		"m":      {"m", "objc", "1.0.0"},
		"ml":     {"ml", "ocaml", "1.0.0"},
		"pl":     {"pl", "perl", "1.0.0"},
		"php":    {"php", "php", "1.0.0"},
		"py":     {"py", "python", "1.0.0"},
		"rkt":    {"rkt", "racket", "1.0.0"},
		"rb":     {"rb", "ruby", "1.0.0"},
		"rs":     {"rs", "rust", "1.0.0"},
		"scala":  {"scala", "scala", "1.0.0"},
		"sh":     {"sh", "bash", "1.0.0"},
	}
	return func(key string) DexecImage {
		return innerMap[key]
	}
}()

const dexecPath = "/tmp/dexec/build"
const dexecImageTemplate = "dexec/%s:%s"
const dexecVolumeTemplate = "%s/%s:%s/%s"

// ExtractBasenameAndPermission takes an include string and splits it into
// its file or folder name and the permission string if present or the empty
// string if not.
func ExtractBasenameAndPermission(path string) (string, string) {
	pathPattern := regexp.MustCompile("([\\w.-]+)(:(rw|ro))")
	match := pathPattern.FindStringSubmatch(path)

	basename := path
	var permission string

	if len(match) == 4 {
		basename = match[1]
		permission = match[2]
	}
	return basename, permission
}

// RunDexecContainer runs an anonymous Docker container with a Docker Exec
// image, mounting the specified sources and includes and passing the
// list of sources and arguments to the entrypoint.
func RunDexecContainer(dexecImage DexecImage, options map[OptionType][]string) {
	dockerImage := fmt.Sprintf(dexecImageTemplate, dexecImage.image, dexecImage.version)

	path := "."
	if len(options[TargetDir]) > 0 {
		path = options[TargetDir][0]
	}
	absPath, _ := filepath.Abs(path)

	var dockerArgs []string
	for _, source := range append(options[Source], options[Include]...) {
		basename, _ := ExtractBasenameAndPermission(source)

		dockerArgs = append(
			dockerArgs,
			[]string{
				"-v",
				fmt.Sprintf(dexecVolumeTemplate, absPath, basename, dexecPath, source),
			}...,
		)
	}

	var sourceBasenames []string
	for _, source := range options[Source] {
		basename, _ := ExtractBasenameAndPermission(source)
		sourceBasenames = append(sourceBasenames, []string{basename}...)
	}

	entrypointArgs := JoinStringSlices(
		sourceBasenames,
		AddPrefix(options[BuildArg], "-b"),
		AddPrefix(options[Arg], "-a"),
	)

	if len(options[UpdateFlag]) > 0 {
		DockerPull(dockerImage)
	}

	RunAnonymousContainer(
		dockerImage,
		dockerArgs,
		entrypointArgs,
	)
}

func validate(cli CLI) bool {
	if !IsDockerPresent() {
		log.Fatal("Docker not found")
	} else if !IsDockerRunning() {
		log.Fatal("Docker not running")
	}

	valid := false
	if len(cli.options[VersionFlag]) != 0 {
		DisplayVersion(cli.filename)
	} else if len(cli.options[Source]) == 0 ||
		len(cli.options[HelpFlag]) != 0 ||
		len(cli.options[TargetDir]) > 1 {
		DisplayHelp(cli.filename)
	} else {
		valid = true
	}
	return valid
}

func main() {
	cli := ParseOsArgs(os.Args)

	if validate(cli) {
		RunDexecContainer(
			LookupImageByExtension(ExtractFileExtension(cli.options[Source][0])),
			cli.options,
		)
	}
}
