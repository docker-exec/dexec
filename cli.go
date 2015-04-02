package main

import (
	"errors"
	"fmt"
	"regexp"
)

type OptionType int

const (
	None     OptionType = iota
	Arg      OptionType = iota
	BuildArg OptionType = iota
	Source   OptionType = iota
)

type Options struct {
	filename string
	options  map[OptionType][]string
}

func GetTypeForOpt(opt string, next string) (OptionType, string, int, error) {
	patternStandaloneA := regexp.MustCompile(`^-(a|-arg)$`)
	patternStandaloneB := regexp.MustCompile(`^-(b|-build-arg)$`)
	patternCombinationA := regexp.MustCompile(`^--arg=(.+)$`)
	patternCombinationB := regexp.MustCompile(`^--build-arg=(.+)$`)
	patternSource := regexp.MustCompile(`^[^-_].+\..+`)

	switch {
	case patternStandaloneA.FindStringIndex(opt) != nil:
		return Arg, next, 2, nil
	case patternStandaloneB.FindStringIndex(opt) != nil:
		return BuildArg, next, 2, nil
	case patternCombinationA.FindStringIndex(opt) != nil:
		return Arg, patternCombinationA.FindStringSubmatch(opt)[1], 1, nil
	case patternCombinationB.FindStringIndex(opt) != nil:
		return BuildArg, patternCombinationB.FindStringSubmatch(opt)[1], 1, nil
	case patternSource.FindStringIndex(opt) != nil:
		return Source, opt, 1, nil
	default:
		return None, "", 0, errors.New(fmt.Sprintf("Unknown option: %s", opt))
	}
}

func ParseOptions(options []string) map[OptionType][]string {
	if len(options) == 0 {
		return map[OptionType][]string{}
	}

	next := ""
	if len(options) > 1 {
		next = options[1]
	}
	t, v, c, _ := GetTypeForOpt(options[0], next)

	if len(options) < c {
		return map[OptionType][]string{}
	}

	m := ParseOptions(options[c:])
	m[t] = append([]string{v}, m[t]...)
	return m
}

func ParseOsArgs(osArgs []string) Options {
	var Options Options

	Options.filename = osArgs[0]
	Options.options = ParseOptions(osArgs[1:])

	return Options
}
