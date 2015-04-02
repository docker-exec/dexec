package main

import (
	"reflect"
	"testing"
)

type OptionData struct {
	first  string
	second string
}

type WantedData struct {
	optionType OptionType
	value      string
	chomped    bool
}

func TestGet(t *testing.T) {
	cases := []struct {
		opt  OptionData
		want WantedData
	}{
		{
			OptionData{"foo.bar", ""},
			WantedData{Source, "foo.bar", false},
		},
		{
			OptionData{"-b", "foo"},
			WantedData{BuildArg, "foo", true},
		},
		{
			OptionData{"--build-arg", "foo"},
			WantedData{BuildArg, "foo", true},
		},
		{
			OptionData{"--build-arg=foo", ""},
			WantedData{BuildArg, "foo", false},
		},
		{
			OptionData{"-a", "foo"},
			WantedData{Arg, "foo", true},
		},
		{
			OptionData{"--arg", "foo"},
			WantedData{Arg, "foo", true},
		},
		{
			OptionData{"--arg=foo", ""},
			WantedData{Arg, "foo", false},
		},
	}
	for _, c := range cases {
		gotOption, gotChomped, _ := GetTypeForOpt(c.opt.first, c.opt.second)
		if gotOption.optionType != c.want.optionType {
			t.Errorf("ParseOsArgsR %q != %q", gotOption.optionType, c.want.optionType)
		} else if gotOption.value != c.want.value {
			t.Errorf("ParseOsArgsR %q != %q", gotOption.value, c.want.value)
		} else if gotChomped != c.want.chomped {
			t.Errorf("ParseOsArgsR %q != %q", gotChomped, c.want.chomped)
		}
	}
}

func TestFilename(t *testing.T) {
	cases := []struct {
		osArgs []string
		want   string
	}{
		{[]string{"foo"}, "foo"},
		{[]string{"bar"}, "bar"},
	}
	for _, c := range cases {
		got := ParseOsArgs(c.osArgs)
		if got.filename != c.want {
			t.Errorf("ParseOsArgsR %q != %q", got.filename, c.want)
		}
	}
}

func TestSources(t *testing.T) {
	cases := []struct {
		osArgs []string
		want   []string
	}{
		{
			[]string{"filename", "foo.cpp", "bar.java", "foo.bar.scala", "bar-foo.groovy"},
			[]string{"foo.cpp", "bar.java", "foo.bar.scala", "bar-foo.groovy"},
		},
	}
	for _, c := range cases {
		got := ParseOsArgs(c.osArgs)
		if !reflect.DeepEqual(got.options[Source], c.want) {
			t.Errorf("ParseOsArgs %q != %q", got.options[Source], c.want)
		}
	}
}

func TestArgs(t *testing.T) {
	cases := []struct {
		osArgs []string
		want   []string
	}{
		{
			[]string{"filename", "-a", "foo", "--arg", "bar", "--arg=foobar"},
			[]string{"foo", "bar", "foobar"},
		},
	}
	for _, c := range cases {
		got := ParseOsArgs(c.osArgs)
		if !reflect.DeepEqual(got.options[Arg], c.want) {
			t.Errorf("ParseOsArgs %q != %q", got.options[Arg], c.want)
		}
	}
}

func TestBuildArgs(t *testing.T) {
	cases := []struct {
		osArgs []string
		want   []string
	}{
		{
			[]string{"filename", "-b", "foo", "--build-arg", "bar", "--build-arg=foobar"},
			[]string{"foo", "bar", "foobar"},
		},
	}
	for _, c := range cases {
		got := ParseOsArgs(c.osArgs)
		if !reflect.DeepEqual(got.options[BuildArg], c.want) {
			t.Errorf("ParseOsArgs %q != %q", got.options[BuildArg], c.want)
		}
	}
}

func TestOrdering(t *testing.T) {
	cases := []struct {
		osArgs []string
		want   map[OptionType][]string
	}{
		{
			[]string{
				"filename",
				"source1.foo",
				"-b", "b_foo",
				"source2.foo",
				"--arg=a_foobar",
				"source3.foo",
				"--build-arg", "b_bar",
				"source4.foo",
				"--arg", "a_bar",
				"source5.foo",
				"--build-arg=b_foobar",
				"source6.foo",
				"-a", "a_foo",
			},
			map[OptionType][]string{
				Arg: {
					"a_foobar", "a_bar", "a_foo",
				},
				BuildArg: {
					"b_foo", "b_bar", "b_foobar",
				},
				Source: {
					"source1.foo",
					"source2.foo",
					"source3.foo",
					"source4.foo",
					"source5.foo",
					"source6.foo",
				},
			},
		},
	}
	for _, c := range cases {
		got := ParseOsArgs(c.osArgs)
		if !reflect.DeepEqual(got.options[BuildArg], c.want[BuildArg]) {
			t.Errorf("ParseOsArgs %q != %q", got.options[BuildArg], c.want[BuildArg])
		} else if !reflect.DeepEqual(got.options[Arg], c.want[Arg]) {
			t.Errorf("ParseOsArgs %q != %q", got.options[Arg], c.want[Arg])
		} else if !reflect.DeepEqual(got.options[Source], c.want[Source]) {
			t.Errorf("ParseOsArgs %q != %q", got.options[Source], c.want[Source])
		}
	}
}
