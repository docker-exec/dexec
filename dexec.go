package main

import (
	"github.com/codegangsta/cli"
	"log"
	"os"
	"regexp"
	"fmt"
	"errors"
)

func GetExtension(filename string) string {
	filenamePattern := regexp.MustCompile(`.*\.(.*)`)
	return filenamePattern.FindStringSubmatch(filename)[1]
}

type argType int

const (
	None argType = iota
	Arg argType = iota
	BuildArg argType = iota
	Source argType = iota
)

type ParsedArgs struct {
	FileName string
	Options map[argType][]string
	Sources, Args, BuildArgs []string
}

func GetTypeForOpt(opt string, next string) (argType, string, bool, error) {
	patternStandaloneA := regexp.MustCompile(`^-(a|-arg)$`)
	patternStandaloneB := regexp.MustCompile(`^-(b|-build-arg)$`)
	patternCombinationA := regexp.MustCompile(`^--arg=(.+)$`)
	patternCombinationB := regexp.MustCompile(`^--build-arg=(.+)$`)
	patternSource := regexp.MustCompile(`^[^-_].+\..+`)

	switch {
		case patternStandaloneA.FindStringIndex(opt) != nil:
			return Arg, next, true, nil
		case patternStandaloneB.FindStringIndex(opt) != nil:
			return BuildArg, next, true, nil
		case patternCombinationA.FindStringIndex(opt) != nil:
			return Arg, patternCombinationA.FindStringSubmatch(opt)[1], false, nil
		case patternCombinationB.FindStringIndex(opt) != nil:
			return BuildArg, patternCombinationB.FindStringSubmatch(opt)[1], false, nil
		case patternSource.FindStringIndex(opt) != nil:
			return Source, opt, false, nil
		default:
			return None, "", false, errors.New("Unknown")
	}
}

func ParseOsArgsRR(osArgs []string) map[argType][]string {
	if len(osArgs) == 0 {
		return map[argType][]string{}
	}

	next := ""
	if len(osArgs) > 1 {
	    next = osArgs[1]
	}
	t, v, c, _ := GetTypeForOpt(osArgs[0], next)

	nextIndex := 1
	if c {
		nextIndex = 2
	}
	if len(osArgs) < nextIndex {
		return map[argType][]string{}
	}

	m := ParseOsArgsRR(osArgs[nextIndex:])
	m[t] = append([]string{v}, m[t]...)
	return m
}

func ParseOsArgs(osArgs []string) ParsedArgs {
	var parsedArgs ParsedArgs

	parsedArgs.FileName = osArgs[0]
	m := ParseOsArgsRR(osArgs[1:])

	parsedArgs.Sources = m[Source]
	parsedArgs.Args = m[Arg]
	parsedArgs.BuildArgs = m[BuildArg]

	return parsedArgs
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
