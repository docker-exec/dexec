package cli

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

	// Image indicates that the option value should be used to
	// override the image worked out based on the file extension.
	Image OptionType = iota

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

	// Extension specifies the override file extenision to use.
	Extension OptionType = iota

	// CleanFlag indicates that the option specifies the clean flag.
	CleanFlag OptionType = iota
)

// CLI defines a data structure that represents the application's name and
// a map of the various options to be used when starting the container.
type CLI struct {
	Filename string
	Options  map[OptionType][]string
}

// ArgToOption takes two candidate strings and returns a tuple consisting of
// what type of option the strings define, the value of the option, how many
// of the strings it took to extract the option value, and a nillable error.
func ArgToOption(opt string, next string) (OptionType, string, int, error) {
	patternStandaloneA := regexp.MustCompile(`^-(a|-arg)$`)
	patternStandaloneB := regexp.MustCompile(`^-(b|-build-arg)$`)
	patternStandaloneI := regexp.MustCompile(`^-(i|-include)$`)
	patternStandaloneM := regexp.MustCompile(`^-(m|-image)$`)
	patternStandaloneE := regexp.MustCompile(`^-(e|-extension)$`)
	patternStandaloneC := regexp.MustCompile(`^-C$`)
	patternCombinationA := regexp.MustCompile(`^--arg=(.+)$`)
	patternCombinationB := regexp.MustCompile(`^--build-arg=(.+)$`)
	patternCombinationI := regexp.MustCompile(`^--include=(.+)$`)
	patternCombinationM := regexp.MustCompile(`^--image=(.+)$`)
	patternCombinationE := regexp.MustCompile(`^--extension=(.+)$`)
	patternSource := regexp.MustCompile(`^[^-_].*\..+`)
	patternUpdateFlag := regexp.MustCompile(`^-(-update|u)$`)
	patternHelpFlag := regexp.MustCompile(`^-(-help|h)$`)
	patternVersionFlag := regexp.MustCompile(`^-(-version|v)$`)
	patternCleanFlag := regexp.MustCompile(`^--clean$`)

	switch {
	case patternStandaloneA.FindStringIndex(opt) != nil:
		return Arg, next, 2, nil
	case patternStandaloneB.FindStringIndex(opt) != nil:
		return BuildArg, next, 2, nil
	case patternStandaloneI.FindStringIndex(opt) != nil:
		return Include, next, 2, nil
	case patternStandaloneM.FindStringIndex(opt) != nil:
		return Image, next, 2, nil
	case patternStandaloneE.FindStringIndex(opt) != nil:
		return Extension, next, 2, nil
	case patternStandaloneC.FindStringIndex(opt) != nil:
		return TargetDir, next, 2, nil
	case patternCombinationA.FindStringIndex(opt) != nil:
		return Arg, patternCombinationA.FindStringSubmatch(opt)[1], 1, nil
	case patternCombinationB.FindStringIndex(opt) != nil:
		return BuildArg, patternCombinationB.FindStringSubmatch(opt)[1], 1, nil
	case patternCombinationI.FindStringIndex(opt) != nil:
		return Include, patternCombinationI.FindStringSubmatch(opt)[1], 1, nil
	case patternCombinationM.FindStringIndex(opt) != nil:
		return Image, patternCombinationM.FindStringSubmatch(opt)[1], 1, nil
	case patternCombinationE.FindStringIndex(opt) != nil:
		return Extension, patternCombinationE.FindStringSubmatch(opt)[1], 1, nil
	case patternUpdateFlag.FindStringIndex(opt) != nil:
		return UpdateFlag, "", 1, nil
	case patternHelpFlag.FindStringIndex(opt) != nil:
		return HelpFlag, "", 1, nil
	case patternVersionFlag.FindStringIndex(opt) != nil:
		return VersionFlag, "", 1, nil
	case patternCleanFlag.FindStringIndex(opt) != nil:
		return CleanFlag, "", 1, nil
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
		Filename: args[0],
		Options:  ParseArgs(args[1:]),
	}
}

// DisplayHelp takes a filename and prints the help information for the program.
func DisplayHelp(filename string) {
	fmt.Println("Name:")
	fmt.Printf("\t%s - Execute code in many languages with Docker!\n", filename)
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Printf("\t%s [options] <source files...>\n", filename)
	fmt.Println()
	fmt.Println("Options:")
	fmt.Printf("\t%-36s%s\n", "-C <dir>", "Specify source directory")
	fmt.Printf("\t%-36s%s\n", "--arg, -a <argument>", "Pass <argument> to the executing code")
	fmt.Printf("\t%-36s%s\n", "--build-arg, -b <build argument>", "Pass <build argument> to compiler")
	fmt.Printf("\t%-36s%s\n", "--include, -i <file|path>", "Mount local <file|path> in dexec container")
	fmt.Printf("\t%-36s%s\n", "--extension, -e <extension>", "Override the image used by <extension>")
	fmt.Printf("\t%-36s%s\n", "--image, -m <name>", "Override the image used by <name>")
	fmt.Printf("\t%-36s%s\n", "--update, -u", "Force update of image")
	fmt.Printf("\t%-36s%s\n", "--clean", "Remove all local dexec images")
	fmt.Printf("\t%-36s%s\n", "--help, -h", "Show help")
	fmt.Printf("\t%-36s%s\n", "--version, -v", "Display version info")
}

// DisplayVersion prints the version information for the program.
func DisplayVersion(filename string) {
	fmt.Printf("%s 1.0.7\n", filename)
}
