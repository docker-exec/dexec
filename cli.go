package main

import (
	"fmt"
	"regexp"
)

// OptionType allows for the enumeration of different CLI option types.
type OptionType int

const (
	// None indicates that no option type is available.
	None OptionType = iota

	// Arg indicates that the option is an argument to be passed to the
	// executing code.
	Arg OptionType = iota

	// BuildArg indicates that the option is an argument to be passed to the
	// compiler if one is used for the target language.
	BuildArg OptionType = iota

	// Source indicates that the option is a source file.
	Source OptionType = iota

	// Include indicates that the option is a file or folder to be mounted
	// in the Docker container without passing it to the compiler or
	// executing code.
	Include OptionType = iota

	// SpecifyImage indicates that the option value should be used to
	// override the image worked out based on the file extension.
	SpecifyImage OptionType = iota

	// TargetDir indicates that the option specifies a custom location (i.e.
	// not the current working directory) to which the sources are
	// relatively specified.
	TargetDir OptionType = iota

	// UpdateFlag indicates that the option specifies that Docker images should
	// be manually updated before being used.
	UpdateFlag OptionType = iota

	// HelpFlag indicates that the option specifies the help flag.
	HelpFlag OptionType = iota

	// VersionFlag indicates that the option specifies the version flag.
	VersionFlag OptionType = iota
)

// CLI defines a data structure that represents the application's name and
// a map of the various options to be used when starting the container.
type CLI struct {
	filename string
	options  map[OptionType][]string
}

// ArgToOption takes two candidate strings and returns a tuple consisting of
// what type of option the strings define, the value of the option, how many
// of the strings it took to extract the option value, and a nillable error.
func ArgToOption(opt string, next string) (OptionType, string, int, error) {
	patternStandaloneA := regexp.MustCompile(`^-(a|-arg)$`)
	patternStandaloneB := regexp.MustCompile(`^-(b|-build-arg)$`)
	patternStandaloneI := regexp.MustCompile(`^-(i|-include)$`)
	patternStandaloneS := regexp.MustCompile(`^-(s|-specify-image)$`)
	patternStandaloneC := regexp.MustCompile(`^-C$`)
	patternCombinationA := regexp.MustCompile(`^--arg=(.+)$`)
	patternCombinationB := regexp.MustCompile(`^--build-arg=(.+)$`)
	patternCombinationI := regexp.MustCompile(`^--include=(.+)$`)
	patternCombinationS := regexp.MustCompile(`^--specify-image=(.+)$`)
	patternSource := regexp.MustCompile(`^[^-_].+\..+`)
	patternUpdateFlag := regexp.MustCompile(`^-(-update|u)$`)
	patternHelpFlag := regexp.MustCompile(`^-(-help|h)$`)
	patternVersionFlag := regexp.MustCompile(`^-(-version|v)$`)

	switch {
	case patternStandaloneA.FindStringIndex(opt) != nil:
		return Arg, next, 2, nil
	case patternStandaloneB.FindStringIndex(opt) != nil:
		return BuildArg, next, 2, nil
	case patternStandaloneI.FindStringIndex(opt) != nil:
		return Include, next, 2, nil
	case patternStandaloneS.FindStringIndex(opt) != nil:
		return SpecifyImage, next, 2, nil
	case patternStandaloneC.FindStringIndex(opt) != nil:
		return TargetDir, next, 2, nil
	case patternCombinationA.FindStringIndex(opt) != nil:
		return Arg, patternCombinationA.FindStringSubmatch(opt)[1], 1, nil
	case patternCombinationB.FindStringIndex(opt) != nil:
		return BuildArg, patternCombinationB.FindStringSubmatch(opt)[1], 1, nil
	case patternCombinationI.FindStringIndex(opt) != nil:
		return Include, patternCombinationI.FindStringSubmatch(opt)[1], 1, nil
	case patternCombinationS.FindStringIndex(opt) != nil:
		return SpecifyImage, patternCombinationS.FindStringSubmatch(opt)[1], 1, nil
	case patternUpdateFlag.FindStringIndex(opt) != nil:
		return UpdateFlag, "", 1, nil
	case patternHelpFlag.FindStringIndex(opt) != nil:
		return HelpFlag, "", 1, nil
	case patternVersionFlag.FindStringIndex(opt) != nil:
		return VersionFlag, "", 1, nil
	case patternSource.FindStringIndex(opt) != nil:
		return Source, opt, 1, nil
	default:
		return None, "", 0, fmt.Errorf("Unknown option: %s", opt)
	}
}

// ParseArgs take a string slice comprised of sources, includes, flags, switches
// and their values and returns a map of these types to a string slice
// of the values of each type of option.
func ParseArgs(args []string) map[OptionType][]string {
	if len(args) == 0 {
		return map[OptionType][]string{}
	}

	var next string
	if len(args) > 1 {
		next = args[1]
	}

	optionType, optionValue, nextIndex, _ := ArgToOption(args[0], next)

	if len(args) < nextIndex || nextIndex == 0 {
		return map[OptionType][]string{}
	}

	optionMap := ParseArgs(args[nextIndex:])
	optionMap[optionType] = append([]string{optionValue}, optionMap[optionType]...)
	return optionMap
}

// ParseOsArgs takes a string slice representing the full arguments passed to
// the program, including the filename and returns a CLI containing the
// filename and map of option types to their values.
func ParseOsArgs(args []string) CLI {
	return CLI{
		filename: args[0],
		options:  ParseArgs(args[1:]),
	}
}

// DisplayHelp takes a filename and prints the help information for the program.
func DisplayHelp(filename string) {
	fmt.Println("Name:")
	fmt.Printf("\t%s - Execute code in many languages with Docker!\n", filename)
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Printf("\t%s [options]\n", filename)
	fmt.Println()
	fmt.Println("Options:")
	fmt.Printf("\t%-50s%s\n", "<source file>", "Execute source file")
	fmt.Printf("\t%-50s%s\n", "-C <dir>", "Specify source directory")
	fmt.Printf("\t%-50s%s\n", "--arg, -a <argument>", "Pass <argument> to the executing code")
	fmt.Printf("\t%-50s%s\n", "--build-arg, -b <build argument>", "Pass <build argument> to compiler")
	fmt.Printf("\t%-50s%s\n", "--include, -i <file|path>", "Mount local <file|path> in dexec container")
	fmt.Printf("\t%-50s%s\n", "--specify-image, -s <docker image>", "Override the image used with <docker image>")
	fmt.Printf("\t%-50s%s\n", "--update, -u", "Update")
	fmt.Printf("\t%-50s%s\n", "--help, -h", "Show help")
	fmt.Printf("\t%-50s%s\n", "--version, -v", "Display version info")
}

// DisplayVersion prints the version information for the program.
func DisplayVersion(filename string) {
	fmt.Printf("%s 1.0.1\n", filename)
}
