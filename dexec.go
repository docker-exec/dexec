package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"code.google.com/p/go-uuid/uuid"

	"github.com/docker-exec/dexec/cli"
	"github.com/docker-exec/dexec/util"
	"github.com/fsouza/go-dockerclient"
)

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

var innerMap = map[string]DexecImage{
	"c":      {"c", "dexec/lang-c", "1.0.2"},
	"clj":    {"clj", "dexec/lang-clojure", "1.0.1"},
	"coffee": {"coffee", "dexec/lang-coffee", "1.0.2"},
	"cpp":    {"cpp", "dexec/lang-cpp", "1.0.2"},
	"cs":     {"cs", "dexec/lang-csharp", "1.0.2"},
	"d":      {"d", "dexec/lang-d", "1.0.1"},
	"erl":    {"erl", "dexec/lang-erlang", "1.0.1"},
	"fs":     {"fs", "dexec/lang-fsharp", "1.0.2"},
	"go":     {"go", "dexec/lang-go", "1.0.1"},
	"groovy": {"groovy", "dexec/lang-groovy", "1.0.1"},
	"hs":     {"hs", "dexec/lang-haskell", "1.0.1"},
	"java":   {"java", "dexec/lang-java", "1.0.2"},
	"lisp":   {"lisp", "dexec/lang-lisp", "1.0.1"},
	"lua":    {"lua", "dexec/lang-lua", "1.0.1"},
	"js":     {"js", "dexec/lang-node", "1.0.2"},
	"nim":    {"nim", "dexec/lang-nim", "1.0.1"},
	"m":      {"m", "dexec/lang-objc", "1.0.1"},
	"ml":     {"ml", "dexec/lang-ocaml", "1.0.1"},
	"p6":     {"p6", "dexec/lang-perl6", "1.0.1"},
	"pl":     {"pl", "dexec/lang-perl", "1.0.2"},
	"php":    {"php", "dexec/lang-php", "1.0.1"},
	"py":     {"py", "dexec/lang-python", "1.0.2"},
	"r":      {"r", "dexec/lang-r", "1.0.1"},
	"rkt":    {"rkt", "dexec/lang-racket", "1.0.1"},
	"rb":     {"rb", "dexec/lang-ruby", "1.0.1"},
	"rs":     {"rs", "dexec/lang-rust", "1.0.1"},
	"scala":  {"scala", "dexec/lang-scala", "1.0.1"},
	"sh":     {"sh", "dexec/lang-bash", "1.0.1"},
}

// LookupImageByExtension is a closure storing a dictionary mapping source
// extensions to the names and versions of Docker Exec images.
var LookupImageByExtension = func() func(string) DexecImage {
	return func(key string) DexecImage {
		return innerMap[key]
	}
}()

var LookupImageByName = func() func(string) DexecImage {
	return func(name string) DexecImage {
		for _, v := range innerMap {
			if v.image == name {
				return v
			}
		}
		panic("no such image")
	}
}()

const dexecPath = "/tmp/dexec/build"
const dexecImageTemplate = "%s:%s"
const dexecVolumeTemplate = "%s/%s:%s/%s"
const dexecSanitisedWindowsPathPattern = "/%s%s"

// ExtractBasenameAndPermission takes an include string and splits it into
// its file or folder name and the permission string if present or the empty
// string if not.
func ExtractBasenameAndPermission(path string) (string, string) {
	pathPattern := regexp.MustCompile("([\\w.:-]+)(:(rw|ro))")
	match := pathPattern.FindStringSubmatch(path)

	basename := path
	var permission string

	if len(match) == 4 {
		basename = match[1]
		permission = match[2]
	}
	return basename, permission
}

// BuildVolumeArgs takes a base path and returns an array of Docker volume
// arguments. The array takes the form {"-v", "/foo:/bar:[rw|ro]", ...} for
// each source or include.
func BuildVolumeArgs(path string, targets []string) []string {
	var volumeArgs []string

	for _, source := range targets {
		basename, _ := ExtractBasenameAndPermission(source)

		volumeArgs = append(
			volumeArgs,
			fmt.Sprintf(dexecVolumeTemplate, path, basename, dexecPath, source),
		)
	}
	return volumeArgs
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

		sanitised = fmt.Sprintf(dexecSanitisedWindowsPathPattern, driveLetter, pathRemainder)
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

// RunDexecContainer runs an anonymous Docker container with a Docker Exec
// image, mounting the specified sources and includes and passing the
// list of sources and arguments to the entrypoint.
func RunDexecContainer(dexecImage DexecImage, options map[cli.OptionType][]string) {
	useStdin := len(options[cli.Source]) == 0
	dockerImage := fmt.Sprintf(dexecImageTemplate, dexecImage.image, dexecImage.version)

	client, err := docker.NewClientFromEnv()

	image, err := client.InspectImage(dockerImage)

	if len(options[cli.UpdateFlag]) > 0 || image == nil {
		client.PullImage(docker.PullImageOptions{
			Repository: dexecImage.image,
			Tag:        dexecImage.version,
		}, docker.AuthConfiguration{})

		image, err = client.InspectImage(dockerImage)
		if err != nil {
			log.Fatal(err)
		} else if image == nil {
			log.Fatal("image was nil")
		}
	}

	if err != nil {
		log.Fatal(err)
	}

	if useStdin {
		stat, _ := os.Stdin.Stat()
		isPipe := (stat.Mode() & os.ModeCharDevice) == 0
		if !isPipe {
			fmt.Fprintln(os.Stderr, "Enter your code. Ctrl-D to exit")
		}
		lines := []string{}
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}
		if !isPipe {
			fmt.Fprintf(os.Stderr, "<Ctrl-D>\n")
		}
		if err := scanner.Err(); err != nil {
			log.Fatal(scanner.Err())
		}
		newfilename := fmt.Sprintf("%s.%s", uuid.NewUUID().String(), dexecImage.extension)

		util.WriteFile(newfilename, []byte(strings.Join(lines, "\n")))
		options[cli.Source] = []string{newfilename}
	}

	defer func() {
		if useStdin {
			util.DeleteFile(options[cli.Source][0])
		}
	}()

	var sourceBasenames []string
	for _, source := range options[cli.Source] {
		basename, _ := ExtractBasenameAndPermission(source)
		sourceBasenames = append(sourceBasenames, []string{basename}...)
	}

	entrypointArgs := util.JoinStringSlices(
		sourceBasenames,
		util.AddPrefix(options[cli.BuildArg], "-b"),
		util.AddPrefix(options[cli.Arg], "-a"),
	)

	container, err := client.CreateContainer(docker.CreateContainerOptions{
		Config: &docker.Config{
			Image: dockerImage,
			Cmd:   entrypointArgs,
		},
	})

	if err != nil {
		log.Fatal(err)
	}

	err = client.StartContainer(container.ID, &docker.HostConfig{
		Binds: BuildVolumeArgs(
			RetrievePath(options[cli.TargetDir]),
			append(options[cli.Source], options[cli.Include]...)),
	})

	if err != nil {
		log.Fatal(err)
	}

	client.AttachToContainer(docker.AttachToContainerOptions{
		Container:    container.ID,
		OutputStream: os.Stdout,
		ErrorStream:  os.Stderr,
		Stream:       true,
		Stdout:       true,
		Stderr:       true,
		Logs:         true,
	})

	if err != nil {
		log.Fatal(err)
	}

	err = client.RemoveContainer(docker.RemoveContainerOptions{
		ID: container.ID,
	})
	if err != nil {
		log.Fatal(err)
	}
}

func validate(cliParser cli.CLI) bool {
	// if !docker.IsDockerPresent() {
	// 	log.Fatal("Docker not found")
	// } else if !docker.IsDockerRunning() {
	// 	log.Fatal("Docker not running")
	// }

	valid := false
	if len(cliParser.Options[cli.VersionFlag]) != 0 {
		cli.DisplayVersion(cliParser.Filename)
	} else if len(cliParser.Options[cli.HelpFlag]) != 0 ||
		len(cliParser.Options[cli.TargetDir]) > 1 ||
		len(cliParser.Options[cli.SpecifyImage]) > 1 {
		cli.DisplayHelp(cliParser.Filename)
	} else if len(cliParser.Options[cli.Source]) == 0 {
		if len(cliParser.Options[cli.Extension]) == 1 ||
			len(cliParser.Options[cli.SpecifyImage]) == 1 {
			valid = true
		} else {
			cli.DisplayHelp(cliParser.Filename)
		}
	} else {
		valid = true
	}
	return valid
}

func main() {
	cliParser := cli.ParseOsArgs(os.Args)

	if validate(cliParser) {
		useStdin := len(cliParser.Options[cli.Source]) == 0

		var image DexecImage
		if useStdin {
			extensionOverride := len(cliParser.Options[cli.Extension]) == 1
			if extensionOverride {
				image = LookupImageByExtension(cliParser.Options[cli.Extension][0])
			} else {
				overrideImage := LookupImageByOverride(cliParser.Options[cli.SpecifyImage][0], "unknown")
				image = LookupImageByName(overrideImage.image)
				image.version = overrideImage.version
			}
		} else {
			extension := util.ExtractFileExtension(cliParser.Options[cli.Source][0])
			image = LookupImageByExtension(extension)
			imageOverride := len(cliParser.Options[cli.SpecifyImage]) == 1
			extensionOverride := len(cliParser.Options[cli.Extension]) == 1
			if extensionOverride {
				image = LookupImageByExtension(cliParser.Options[cli.Extension][0])
			} else if imageOverride {
				image = LookupImageByOverride(cliParser.Options[cli.SpecifyImage][0], extension)
			}
		}

		RunDexecContainer(
			image,
			cliParser.Options,
		)
	}
}
