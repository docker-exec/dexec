package main

import (
	"errors"
	"fmt"
	"regexp"
)

type OptionType int

const (
	None        OptionType = iota
	Arg         OptionType = iota
	BuildArg    OptionType = iota
	Source      OptionType = iota
	Include     OptionType = iota
	TargetDir   OptionType = iota
	HelpFlag    OptionType = iota
	VersionFlag OptionType = iota
)

type CLI struct {
	filename string
	options  map[OptionType][]string
}

func ArgToOption(opt string, next string) (OptionType, string, int, error) {
	patternStandaloneA := regexp.MustCompile(`^-(a|-arg)$`)
	patternStandaloneB := regexp.MustCompile(`^-(b|-build-arg)$`)
	patternStandaloneI := regexp.MustCompile(`^-(i|-include)$`)
	patternStandaloneC := regexp.MustCompile(`^-C$`)
	patternCombinationA := regexp.MustCompile(`^--arg=(.+)$`)
	patternCombinationB := regexp.MustCompile(`^--build-arg=(.+)$`)
	patternCombinationI := regexp.MustCompile(`^--include=(.+)$`)
	patternSource := regexp.MustCompile(`^[^-_].+\..+`)
	patternHelpFlag := regexp.MustCompile(`^-(-help|h)$`)
	patternVersionFlag := regexp.MustCompile(`^-(-version|v)$`)

	switch {
	case patternStandaloneA.FindStringIndex(opt) != nil:
		return Arg, next, 2, nil
	case patternStandaloneB.FindStringIndex(opt) != nil:
		return BuildArg, next, 2, nil
	case patternStandaloneI.FindStringIndex(opt) != nil:
		return Include, next, 2, nil
	case patternStandaloneC.FindStringIndex(opt) != nil:
		return TargetDir, next, 2, nil
	case patternCombinationA.FindStringIndex(opt) != nil:
		return Arg, patternCombinationA.FindStringSubmatch(opt)[1], 1, nil
	case patternCombinationB.FindStringIndex(opt) != nil:
		return BuildArg, patternCombinationB.FindStringSubmatch(opt)[1], 1, nil
	case patternCombinationI.FindStringIndex(opt) != nil:
		return Include, patternCombinationI.FindStringSubmatch(opt)[1], 1, nil
	case patternHelpFlag.FindStringIndex(opt) != nil:
		return HelpFlag, "", 1, nil
	case patternVersionFlag.FindStringIndex(opt) != nil:
		return VersionFlag, "", 1, nil
	case patternSource.FindStringIndex(opt) != nil:
		return Source, opt, 1, nil
	default:
		return None, "", 0, errors.New(fmt.Sprintf("Unknown option: %s", opt))
	}
}

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

func ParseOsArgs(args []string) CLI {
	return CLI{
		filename: args[0],
		options:  ParseArgs(args[1:]),
	}
}

func DisplayHelp(filename string) {
	fmt.Println("Name:")
	fmt.Printf("\t%s - Execute code in many languages with Docker!\n", filename)
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("\tdexec [options]")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Printf("\t%-50s%s\n", "<source file>", "Execute source file")
	fmt.Printf("\t%-50s%s\n", "-C <dir>", "Specify source directory")
	fmt.Printf("\t%-50s%s\n", "--arg, -a <argument>", "Pass <argument> to the executing code")
	fmt.Printf("\t%-50s%s\n", "--build-arg, -b <build argument>", "Pass <build argument> to compiler")
	fmt.Printf("\t%-50s%s\n", "--help, -h", "Show help")
	fmt.Printf("\t%-50s%s\n", "--version, -v", "Display version info")
}

func DisplayVersion() {
	fmt.Println("dexec 1.0.0-alpha")
}
