package main

import (
	"errors"
	"regexp"
)

type OptionType int

const (
	None     OptionType = iota
	Arg      OptionType = iota
	BuildArg OptionType = iota
	Source   OptionType = iota
)

type Option struct {
	optionType OptionType
	value      string
}

type ParsedArgs struct {
	filename string
	options  map[OptionType][]string
}

func GetTypeForOpt(opt string, next string) (Option, bool, error) {
	patternStandaloneA := regexp.MustCompile(`^-(a|-arg)$`)
	patternStandaloneB := regexp.MustCompile(`^-(b|-build-arg)$`)
	patternCombinationA := regexp.MustCompile(`^--arg=(.+)$`)
	patternCombinationB := regexp.MustCompile(`^--build-arg=(.+)$`)
	patternSource := regexp.MustCompile(`^[^-_].+\..+`)

	switch {
	case patternStandaloneA.FindStringIndex(opt) != nil:
		return Option{Arg, next}, true, nil
	case patternStandaloneB.FindStringIndex(opt) != nil:
		return Option{BuildArg, next}, true, nil
	case patternCombinationA.FindStringIndex(opt) != nil:
		return Option{Arg, patternCombinationA.FindStringSubmatch(opt)[1]}, false, nil
	case patternCombinationB.FindStringIndex(opt) != nil:
		return Option{BuildArg, patternCombinationB.FindStringSubmatch(opt)[1]}, false, nil
	case patternSource.FindStringIndex(opt) != nil:
		return Option{Source, opt}, false, nil
	default:
		return Option{None, ""}, false, errors.New("Unknown")
	}
}

func ParseOsArgsRR(osArgs []string) map[OptionType][]string {
	if len(osArgs) == 0 {
		return map[OptionType][]string{}
	}

	next := ""
	if len(osArgs) > 1 {
		next = osArgs[1]
	}
	o, c, _ := GetTypeForOpt(osArgs[0], next)

	nextIndex := 1
	if c {
		nextIndex = 2
	}
	if len(osArgs) < nextIndex {
		return map[OptionType][]string{}
	}

	m := ParseOsArgsRR(osArgs[nextIndex:])
	m[o.optionType] = append([]string{o.value}, m[o.optionType]...)
	return m
}

func ParseOsArgs(osArgs []string) ParsedArgs {
	var parsedArgs ParsedArgs

	parsedArgs.filename = osArgs[0]
	parsedArgs.options = ParseOsArgsRR(osArgs[1:])

	return parsedArgs
}
