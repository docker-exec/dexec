package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"code.google.com/p/go-uuid/uuid"

	"github.com/docker-exec/dexec/cli"
	"github.com/docker-exec/dexec/util"
	"github.com/fsouza/go-dockerclient"
)

// DexecImage consists of the file extension, Docker image name and Docker
// image version to use for a given Docker Exec image.
type DexecImage struct {
	name      string
	extension string
	image     string
	version   string
}

const dexecPath = "/tmp/dexec/build"
const dexecImageTemplate = "%s:%s"
const dexecVolumeTemplate = "%s/%s:%s/%s"

var innerMap = map[string]*DexecImage{
	"c":      {"C", "c", "dexec/lang-c", "1.0.2"},
	"clj":    {"Clojure", "clj", "dexec/lang-clojure", "1.0.1"},
	"coffee": {"CoffeeScript", "coffee", "dexec/lang-coffee", "1.0.2"},
	"cpp":    {"C++", "cpp", "dexec/lang-cpp", "1.0.2"},
	"cs":     {"C#", "cs", "dexec/lang-csharp", "1.0.2"},
	"d":      {"D", "d", "dexec/lang-d", "1.0.1"},
	"erl":    {"Erlang", "erl", "dexec/lang-erlang", "1.0.1"},
	"fs":     {"F#", "fs", "dexec/lang-fsharp", "1.0.2"},
	"go":     {"Go", "go", "dexec/lang-go", "1.0.1"},
	"groovy": {"Groovy", "groovy", "dexec/lang-groovy", "1.0.1"},
	"hs":     {"Haskell", "hs", "dexec/lang-haskell", "1.0.1"},
	"java":   {"Java", "java", "dexec/lang-java", "1.0.2"},
	"lisp":   {"Lisp", "lisp", "dexec/lang-lisp", "1.0.1"},
	"lua":    {"Lua", "lua", "dexec/lang-lua", "1.0.1"},
	"js":     {"JavaScript", "js", "dexec/lang-node", "1.0.2"},
	"nim":    {"Nim", "nim", "dexec/lang-nim", "1.0.1"},
	"m":      {"Objective C", "m", "dexec/lang-objc", "1.0.1"},
	"ml":     {"OCaml", "ml", "dexec/lang-ocaml", "1.0.1"},
	"p6":     {"Perl 6", "p6", "dexec/lang-perl6", "1.0.1"},
	"pl":     {"Perl", "pl", "dexec/lang-perl", "1.0.2"},
	"php":    {"PHP", "php", "dexec/lang-php", "1.0.1"},
	"py":     {"Python", "py", "dexec/lang-python", "1.0.2"},
	"r":      {"R", "r", "dexec/lang-r", "1.0.1"},
	"rkt":    {"Racket", "rkt", "dexec/lang-racket", "1.0.1"},
	"rb":     {"Ruby", "rb", "dexec/lang-ruby", "1.0.1"},
	"rs":     {"Rust", "rs", "dexec/lang-rust", "1.0.1"},
	"scala":  {"Scala", "scala", "dexec/lang-scala", "1.0.1"},
	"sh":     {"Bash", "sh", "dexec/lang-bash", "1.0.1"},
}

// LookupImageByExtension returns the image for a given extension.
func LookupImageByExtension(key string) (*DexecImage, error) {
	if val, ok := innerMap[key]; ok {
		return val, nil
	}
	return nil, fmt.Errorf("Map does not contain key")
}

// LookupImageByName returns the image for a given image name.
func LookupImageByName(name string) (*DexecImage, error) {
	for _, v := range innerMap {
		if v.image == name {
			return v, nil
		}
	}
	return nil, fmt.Errorf("Map does not contain key")
}

// LookupImageByOverride takes an image that has been specified by the user
// to use instead of the one in the extension map. This function returns a
// DexecImage struct containing the image name & version, as well as the
// file extension that was passed in.
func LookupImageByOverride(image string, extension string) (*DexecImage, error) {
	patternImage := regexp.MustCompile(`(.*):(.*)`)
	imageMatch := patternImage.FindStringSubmatch(image)
	if len(imageMatch) > 0 {
		return &DexecImage{
			"Unknown",
			extension,
			imageMatch[1],
			imageMatch[2],
		}, nil
	}
	return &DexecImage{
		"Unknown",
		extension,
		image,
		"latest",
	}, nil
}

func imageFromOptions(options map[cli.OptionType][]string) *DexecImage {
	var image *DexecImage

	if useStdin := len(options[cli.Source]) == 0; useStdin {
		extensionOverride := len(options[cli.Extension]) == 1
		if extensionOverride {
			image, _ = LookupImageByExtension(options[cli.Extension][0])
		} else {
			overrideImage, _ := LookupImageByOverride(options[cli.SpecifyImage][0], "unknown")
			image, _ = LookupImageByName(overrideImage.image)
			image.version = overrideImage.version
		}
	} else {
		extension := util.ExtractFileExtension(options[cli.Source][0])
		image, _ = LookupImageByExtension(extension)
		imageOverride := len(options[cli.SpecifyImage]) == 1
		extensionOverride := len(options[cli.Extension]) == 1
		if extensionOverride {
			image, _ = LookupImageByExtension(options[cli.Extension][0])
		} else if imageOverride {
			image, _ = LookupImageByOverride(options[cli.SpecifyImage][0], extension)
		}
	}
	return image
}

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

func fetchImage(name string, tag string, update bool, client *docker.Client) error {
	dockerImage := fmt.Sprintf(dexecImageTemplate, name, tag)

	if _, err := client.InspectImage(dockerImage); update || err != nil {
		client.PullImage(docker.PullImageOptions{
			Repository: name,
			Tag:        tag,
		}, docker.AuthConfiguration{})

		if _, err = client.InspectImage(dockerImage); err != nil {
			return err
		}
	}
	return nil
}

// RunDexecContainer runs an anonymous Docker container with a Docker Exec
// image, mounting the specified sources and includes and passing the
// list of sources and arguments to the entrypoint.
func RunDexecContainer(options map[cli.OptionType][]string) {
	dexecImage := imageFromOptions(options)
	dockerImage := fmt.Sprintf(dexecImageTemplate, dexecImage.image, dexecImage.version)
	updateImage := len(options[cli.UpdateFlag]) > 0

	client, err := docker.NewClientFromEnv()

	if err != nil {
		log.Fatal(err)
	}

	if err := fetchImage(
		dexecImage.image,
		dexecImage.version,
		updateImage,
		client); err != nil {
		log.Fatal(err)
	}

	if useStdin := len(options[cli.Source]) == 0; useStdin {
		lines := util.ReadStdin("Enter your code. Ctrl-D to exit", "<Ctrl-D>")
		tmpFile := fmt.Sprintf("%s.%s", uuid.NewUUID().String(), dexecImage.extension)

		util.WriteFile(tmpFile, []byte(strings.Join(lines, "\n")))
		defer func() {
			util.DeleteFile(tmpFile)
		}()

		options[cli.Source] = []string{tmpFile}
	}

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

	if err = client.StartContainer(container.ID, &docker.HostConfig{
		Binds: BuildVolumeArgs(
			util.RetrievePath(options[cli.TargetDir]),
			append(options[cli.Source], options[cli.Include]...)),
	}); err != nil {
		log.Fatal(err)
	}

	if err = client.AttachToContainer(docker.AttachToContainerOptions{
		Container:    container.ID,
		OutputStream: os.Stdout,
		ErrorStream:  os.Stderr,
		Stream:       true,
		Stdout:       true,
		Stderr:       true,
		Logs:         true,
	}); err != nil {
		log.Fatal(err)
	}

	if err = client.RemoveContainer(docker.RemoveContainerOptions{
		ID: container.ID,
	}); err != nil {
		log.Fatal(err)
	}
}

func validate(cliParser cli.CLI) bool {
	options := cliParser.Options

	hasVersionFlag := len(options[cli.VersionFlag]) == 1
	hasExtension := len(options[cli.Extension]) == 1
	hasImage := len(options[cli.SpecifyImage]) == 1
	hasSources := len(options[cli.Source]) > 0

	if hasSources || hasImage || hasExtension {
		return true
	}

	if hasVersionFlag {
		cli.DisplayVersion(cliParser.Filename)
		return false
	}

	cli.DisplayHelp(cliParser.Filename)
	return false
}

func validateDocker() error {
	client, err := docker.NewClientFromEnv()
	if err != nil {
		return err
	}

	ping := make(chan error, 1)
	go func() {
		ping <- client.Ping()
	}()

	select {
	case err := <-ping:
		return err
	case <-time.After(5 * time.Second):
		return fmt.Errorf("Request to Docker host timed out")
	}
}

func main() {
	cliParser := cli.ParseOsArgs(os.Args)

	if validate(cliParser) {
		if err := validateDocker(); err != nil {
			log.Fatal(err)
		} else {
			RunDexecContainer(cliParser.Options)
		}
	}
}
