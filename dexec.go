package main

import (
	"github.com/codegangsta/cli"
	"log"
	"os"
	"regexp"
)

func GetExtension(filename string) string {
	filenamePattern := regexp.MustCompile(`.*\.(.*)`)
	return filenamePattern.FindStringSubmatch(string(filename))[1]
}

func main() {
	extensionMap := map[string]string{
		"c": "c",
		"clj": "clojure",
		"coffee": "coffee",
		"cpp": "cpp",
		"cs": "csharp",
		"d": "d",
		"erl": "erlang",
		"fs": "fsharp",
		"go": "go",
		"groovy": "groovy",
		"hs": "haskell",
		"java": "java",
		"lisp": "lisp",
		"js": "node",
		"m": "objc",
		"ml": "ocaml",
		"pl": "perl",
		"php": "php",
		"py": "python",
		"rkt": "racket",
		"rb": "ruby",
		"rs": "rust",
		"scala": "scala",
		"sh": "bash",
	}

	app := cli.NewApp()
	app.Name = "dexec"
	app.Usage = "dexec"

	app.Action = func(c *cli.Context) {
		found := IsDockerPresent()
		running := IsDockerRunning()

		if !found {
			log.Fatal("Docker not found")
		} else if !running {
			log.Fatal("Docker not running")
		} else if len(c.Args()) == 0 {
			cli.ShowAppHelp(c)
		} else {
			sourceFile := c.Args()[0]
			imageName := extensionMap[GetExtension(sourceFile)]
			RunDexecContainer(imageName, sourceFile, c.Args()[1:]...)
		}
	}

	app.Run(os.Args)
}
