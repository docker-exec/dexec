package main

import (
	"log"
	"os"
	"regexp"
)

func GetExtension(filename string) string {
	filenamePattern := regexp.MustCompile(`.*\.(.*)`)
	return filenamePattern.FindStringSubmatch(filename)[1]
}

var extensionDict = func() func(string) string {
	innerMap := map[string]string{
		"c":      "c",
		"clj":    "clojure",
		"coffee": "coffee",
		"cpp":    "cpp",
		"cs":     "csharp",
		"d":      "d",
		"erl":    "erlang",
		"fs":     "fsharp",
		"go":     "go",
		"groovy": "groovy",
		"hs":     "haskell",
		"java":   "java",
		"lisp":   "lisp",
		"js":     "node",
		"m":      "objc",
		"ml":     "ocaml",
		"pl":     "perl",
		"php":    "php",
		"py":     "python",
		"rkt":    "racket",
		"rb":     "ruby",
		"rs":     "rust",
		"scala":  "scala",
		"sh":     "bash",
	}
	return func(key string) string {
		return innerMap[key]
	}
}()


func main() {
	found := IsDockerPresent()
	running := IsDockerRunning()

	if !found {
		log.Fatal("Docker not found")
	} else if !running {
		log.Fatal("Docker not running")
	} else {
		options := ParseOsArgs(os.Args)

		if len(options.options[VersionFlag]) != 0 {
			PrintVersion()
		} else if len(options.options[Source]) == 0 ||
			len(options.options[HelpFlag]) != 0 {
			PrintHelp()
		} else {
			imageName := extensionDict(GetExtension(options.options[Source][0]))
			RunDexecContainer(imageName, options.options)
		}
	}
}
