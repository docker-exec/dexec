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

// LookupImageByOverride takes an image that has been specified by the user
// to use instead of the one in the extension map. This function returns a
// DexecImage struct containing the image name & version, as well as the
// file extension that was passed in.
func LookupImageByOverride(image string, extension string) DexecImage {
	patternImage := regexp.MustCompile(`(.*):(.*)`)
	imageMatch := patternImage.FindStringSubmatch(image)
	if len(imageMatch) > 0 {
		return DexecImage{
			extension,
			imageMatch[1],
			imageMatch[2],
		}
	}
	return DexecImage{
		extension,
		image,
		"latest",
	}
}

// LookupImageByExtension is a closure storing a dictionary mapping source
// extensions to the names and versions of Docker Exec images.
var LookupImageByExtension = func() func(string) DexecImage {
	innerMap := map[string]DexecImage{
		"c":      {"c", "dexec/c", "1.0.0"},
		"clj":    {"clj", "dexec/clojure", "1.0.0"},
		"coffee": {"coffee", "dexec/coffee", "1.0.0"},
		"cpp":    {"cpp", "dexec/cpp", "1.0.0"},
		"cs":     {"cs", "dexec/csharp", "1.0.0"},
		"d":      {"d", "dexec/d", "1.0.0"},
		"erl":    {"erl", "dexec/erlang", "1.0.0"},
		"fs":     {"fs", "dexec/fsharp", "1.0.0"},
		"go":     {"go", "dexec/go", "1.0.0"},
		"groovy": {"groovy", "dexec/groovy", "1.0.0"},
		"hs":     {"hs", "dexec/haskell", "1.0.0"},
		"java":   {"java", "dexec/java", "1.0.0"},
		"lisp":   {"lisp", "dexec/lisp", "1.0.0"},
		"lua":    {"lua", "dexec/lua", "latest"},
		"js":     {"js", "dexec/node", "1.0.0"},
		"nim":    {"nim", "dexec/nim", "latest"},
		"m":      {"m", "dexec/objc", "1.0.0"},
		"ml":     {"ml", "dexec/ocaml", "1.0.0"},
		"p6":     {"p6", "dexec/perl6", "1.0.0"},
		"pl":     {"pl", "dexec/perl", "1.0.1"},
		"php":    {"php", "dexec/php", "1.0.0"},
		"py":     {"py", "dexec/python", "1.0.0"},
		"r":      {"r", "dexec/r", "latest"},
		"rkt":    {"rkt", "dexec/racket", "1.0.0"},
		"rb":     {"rb", "dexec/ruby", "1.0.0"},
		"rs":     {"rs", "dexec/rust", "1.0.0"},
		"scala":  {"scala", "dexec/scala", "1.0.0"},
		"sh":     {"sh", "dexec/bash", "1.0.0"},
	}
	return func(key string) DexecImage {
		return innerMap[key]
	}
}()

const dexecPath = "/tmp/dexec/build"
const dexecImageTemplate = "%s:%s"
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
		len(cli.options[TargetDir]) > 1 ||
		len(cli.options[SpecifyImage]) > 1 {
		DisplayHelp(cli.filename)
	} else {
		valid = true
	}
	return valid
}

func main() {
	cli := ParseOsArgs(os.Args)

	if validate(cli) {
		extension := ExtractFileExtension(cli.options[Source][0])
		image := LookupImageByExtension(extension)
		if len(cli.options[SpecifyImage]) == 1 {
			image = LookupImageByOverride(cli.options[SpecifyImage][0], extension)
		}
		RunDexecContainer(
			image,
			cli.options,
		)
	}
}
