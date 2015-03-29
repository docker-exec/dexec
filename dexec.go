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
	return filenamePattern.FindStringSubmatch(filename)[1]
}

type ArgType int

const (
	None ArgType = iota
	Arg ArgType = iota
	BuildArg ArgType = iota
	Source ArgType = iota
)


type ParsedArgs struct {
	FileName string
	Sources, Args, BuildArgs []string
}

func ParseOsArgsR(osArgs []string) ParsedArgs {
	var parsedArgs ParsedArgs
	parsedArgs.FileName = osArgs[0]

	parsedArgs.Args, parsedArgs.BuildArgs, parsedArgs.Sources = ParseOsArgsR2(None, osArgs[1:])

	return parsedArgs
}

func ParseOsArgsR2(argType ArgType, osArgs []string) ([]string, []string, []string) {
	if len(osArgs) == 1 {

	}


	switch argType {
		case None: s1()
		case 4, 5, 6, 7: s2()
	}
	return []string{}, []string{}, []string{}
}

func GetArgType(arg string) ArgType {
	
}

func ParseOsArgs(osArgs []string) ParsedArgs {

	patternStandaloneA := regexp.MustCompile(`^-(a|-arg)$`)
	patternStandaloneB := regexp.MustCompile(`^-(b|-build-arg)$`)
	patternCombinationA := regexp.MustCompile(`^--arg=(.+)$`)
	patternCombinationB := regexp.MustCompile(`^--build-arg=(.+)$`)

	patternSource := regexp.MustCompile(`^[^-_].+\..+`)

	var parsedArgs ParsedArgs

	parsedArgs.FileName = osArgs[0]

	toParse := osArgs[1:]

	for i, opt := range toParse {
		if patternStandaloneA.FindStringIndex(opt) != nil {
			parsedArgs.Args = append(parsedArgs.Args, toParse[i + 1])
		} else if patternCombinationA.FindStringIndex(opt) != nil {
			parsedArgs.Args = append(parsedArgs.Args, patternCombinationA.FindStringSubmatch(opt)[1])
		} else if patternStandaloneB.FindStringIndex(opt) != nil {
			parsedArgs.BuildArgs = append(parsedArgs.BuildArgs, toParse[i + 1])
		} else if patternCombinationB.FindStringIndex(opt) != nil {
			parsedArgs.BuildArgs = append(parsedArgs.BuildArgs, patternCombinationB.FindStringSubmatch(opt)[1])
		} else if patternSource.FindStringIndex(opt) != nil {
			parsedArgs.Sources = append(parsedArgs.Sources, opt)
		}
	}

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
