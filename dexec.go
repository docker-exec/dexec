package main

import (
	"fmt"
	"log"
	"regexp"

	"github.com/fsouza/go-dockerclient"
)

// ContainerImage consists of the file extension, Docker image name and Docker
// image version to use for a given Docker Exec image.
type ContainerImage struct {
	Name      string
	Extension string
	Image     string
	Version   string
}

const dexecPath = "/tmp/dexec/build"
const dexecImageTemplate = "%s:%s"
const dexecVolumeTemplate = "%s/%s:%s/%s"

// ImageFromOptions returns an image from a set of options.
func ImageFromOptions(options map[OptionType][]string) (image *ContainerImage, err error) {
	useExtension := len(options[Extension]) == 1
	useImage := len(options[Image]) == 1

	if useStdin := len(options[Source]) == 0; useStdin {
		if useExtension {
			image, err = LookupImageByExtension(options[Extension][0])
		} else if useImage {
			overrideImage, err := LookupImageByOverride(options[Image][0], "unknown")
			if err != nil {
				return nil, err
			}
			image, err = LookupImageByName(overrideImage.Image)
			if image != nil {
				image.Version = overrideImage.Version
			}
		} else {
			err = fmt.Errorf("STDIN requested but no extension or image supplied")
		}
	} else {
		if extension := ExtractFileExtension(options[Source][0]); useExtension {
			image, err = LookupImageByExtension(options[Extension][0])
		} else if useImage {
			image, err = LookupImageByOverride(options[Image][0], extension)
		} else {
			image, err = LookupImageByExtension(extension)
		}
	}
	return image, err
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

// FetchImage guarantees a Docker image is availabe in the local repository or
// returns an error.
func FetchImage(name string, tag string, update bool, client *docker.Client) error {
	dockerImage := fmt.Sprintf(dexecImageTemplate, name, tag)

	if _, err := client.InspectImage(dockerImage); update || err != nil {
		err = client.PullImage(docker.PullImageOptions{
			Repository: name,
			Tag:        tag,
		}, docker.AuthConfiguration{})

		if err != nil {
			log.Fatal(err)
		}

		if _, err = client.InspectImage(dockerImage); err != nil {
			return err
		}
	}
	return nil
}

var innerMap = map[string]*ContainerImage{
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
	"java":   {"Java", "java", "dexec/lang-java", "1.0.3"},
	"lisp":   {"Lisp", "lisp", "dexec/lang-lisp", "1.0.1"},
	"lua":    {"Lua", "lua", "dexec/lang-lua", "1.0.1"},
	"js":     {"JavaScript", "js", "dexec/lang-node", "1.0.2"},
	"nim":    {"Nim", "nim", "dexec/lang-nim", "1.0.1"},
	"m":      {"Objective C", "m", "dexec/lang-objc", "1.0.2"},
	"ml":     {"OCaml", "ml", "dexec/lang-ocaml", "1.0.1"},
	"p6":     {"Perl 6", "p6", "dexec/lang-perl6", "1.0.1"},
	"pl":     {"Perl", "pl", "dexec/lang-perl", "1.0.2"},
	"php":    {"PHP", "php", "dexec/lang-php", "1.0.1"},
	"py":     {"Python", "py", "dexec/lang-python", "1.0.2"},
	"r":      {"R", "r", "dexec/lang-r", "1.0.1"},
	"rkt":    {"Racket", "rkt", "dexec/lang-racket", "1.0.1"},
	"rb":     {"Ruby", "rb", "dexec/lang-ruby", "1.0.2"},
	"rs":     {"Rust", "rs", "dexec/lang-rust", "1.0.1"},
	"scala":  {"Scala", "scala", "dexec/lang-scala", "1.0.1"},
	"sh":     {"Bash", "sh", "dexec/lang-bash", "1.0.1"},
}

// LookupImageByExtension returns the image for a given extension.
func LookupImageByExtension(key string) (*ContainerImage, error) {
	if v, ok := innerMap[key]; ok {
		return v, nil
	}
	return nil, fmt.Errorf("map does not contain key %s", key)
}

// LookupImageByName returns the image for a given image name.
func LookupImageByName(name string) (*ContainerImage, error) {
	for _, v := range innerMap {
		if v.Image == name {
			return v, nil
		}
	}
	return nil, fmt.Errorf("map does not contain image with name %s", name)
}

// LookupImageByOverride takes an image that has been specified by the user
// to use instead of the one in the extension map. This function returns a
// DexecImage struct containing the image name & version, as well as the
// file extension that was passed in.
func LookupImageByOverride(image string, extension string) (*ContainerImage, error) {
	patternImage := regexp.MustCompile(`(.*):(.*)`)
	imageMatch := patternImage.FindStringSubmatch(image)
	if len(imageMatch) > 0 {
		return &ContainerImage{
			"Unknown",
			extension,
			imageMatch[1],
			imageMatch[2],
		}, nil
	}
	return &ContainerImage{
		"Unknown",
		extension,
		image,
		"latest",
	}, nil
}
