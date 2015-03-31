package main

import (
	"github.com/codegangsta/cli"
	"log"
	"os"
	"regexp"
	"fmt"
)

func GetExtension(filename string) string {
	filenamePattern := regexp.MustCompile(`.*\.(.*)`)
	return filenamePattern.FindStringSubmatch(string(filename))[1]
}

func main() {
	extensionMap := map[string]string{
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

	var AppHelpTemplate = `Name:
    {{.Name}} - {{.Usage}}

Usage:
   {{.Name}} [options] [sources]

Options:
   {{range .Flags}}{{.}}
   {{end}}
`

	cli.AppHelpTemplate = AppHelpTemplate

	app := cli.NewApp()
	app.Name = "dexec"
	app.Usage = "Execute code in many languages with Docker!"
	app.Version = "1.0.0-beta"
	app.Author = "Andy Stanton"
	app.EnableBashCompletion = true

	argFlags := cli.StringSlice{}
	buildArgFlags := cli.StringSlice{}

	app.Flags = []cli.Flag{
		cli.StringSliceFlag{
			Name:  "arg, a",
			Usage: "Arguments to pass to the program",
			Value: &argFlags,
		},
		cli.StringSliceFlag{
			Name:  "build-arg, b",
			Usage: "Arguments to pass to the compiler (if the target language has one)",
			Value: &buildArgFlags,
		},
	}

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

			sources := c.Args()
			fmt.Printf("%v\n", sources)

			buildArgs := c.StringSlice("build-arg")
			fmt.Printf("%v\n", buildArgs)

			// imageName := extensionMap[GetExtension(sourceFile)]
			imageName := extensionMap[GetExtension("blah.cpp")]
			fmt.Printf(imageName)
			fmt.Printf(sourceFile)
			// RunDexecContainer(imageName, sourceFile, c.Args()[1:]...)
		}
	}

	app.Run(os.Args)
}
