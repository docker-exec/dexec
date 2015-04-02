package main

import (
	"testing"
	"reflect"
)

type OptStruct struct {
	first string
	second string
}

type WantStruct struct {
	ArgType argType
	Chomped bool
	Value string
}

func TestGet(t *testing.T) {
	cases := []struct {
		opt OptStruct
		want WantStruct
	}{
		{
			OptStruct{"foo.bar", ""},
			WantStruct{Source, false, "foo.bar"},
		},
		{
			OptStruct{"-b", "foo"},
			WantStruct{BuildArg, true, "foo"},
		},
		{
			OptStruct{"--build-arg", "foo"},
			WantStruct{BuildArg, true, "foo"},
		},
		{
			OptStruct{"--build-arg=foo", ""},
			WantStruct{BuildArg, false, "foo"},
		},
		{
			OptStruct{"-a", "foo"},
			WantStruct{Arg, true, "foo"},
		},
		{
			OptStruct{"--arg", "foo"},
			WantStruct{Arg, true, "foo"},
		},
		{
			OptStruct{"--arg=foo", ""},
			WantStruct{Arg, false, "foo"},
		},
	}
	for _, c := range cases {
		gotType, gotValue, gotChomped, _ := GetTypeForOpt(c.opt.first, c.opt.second)
		if gotType != c.want.ArgType {
			t.Errorf("ParseOsArgsR %q != %q", gotType, c.want.ArgType)
		} else if gotChomped != c.want.Chomped {
			t.Errorf("ParseOsArgsR %q != %q", gotChomped, c.want.Chomped)
		} else if gotValue != c.want.Value {
			t.Errorf("ParseOsArgsR %q != %q", gotValue, c.want.Value)
		}
	}
}

func TestFilename(t *testing.T) {
	cases := []struct {
		osArgs []string
		want string
	}{
		{ []string{"foo"}, "foo" },
		{ []string{"bar"}, "bar" },
	}
	for _, c := range cases {
		got := ParseOsArgs(c.osArgs)
		if got.FileName != c.want {
			t.Errorf("ParseOsArgsR %q != %q", got.FileName, c.want)
		}
	}
}

func TestSources(t *testing.T) {
	cases := []struct {
		osArgs []string
		want []string
	}{
		{
			[]string{"Filename", "foo.cpp", "bar.java", "foo.bar.scala", "bar-foo.groovy"},
			[]string{"foo.cpp", "bar.java", "foo.bar.scala", "bar-foo.groovy"},
		},
	}
	for _, c := range cases {
		got := ParseOsArgs(c.osArgs)
		if ! reflect.DeepEqual(got.Sources, c.want) {
			t.Errorf("ParseOsArgs %q != %q", got.Sources, c.want)
		}
	}
}

func TestArgs(t *testing.T) {
	cases := []struct {
		osArgs []string
		want []string
	}{
		{
			[]string{"Filename", "-a", "foo", "--arg", "bar", "--arg=foobar"},
			[]string{"foo", "bar", "foobar"},
		},
	}
	for _, c := range cases {
		got := ParseOsArgs(c.osArgs)
		if ! reflect.DeepEqual(got.Args, c.want) {
			t.Errorf("ParseOsArgs %q != %q", got.Args, c.want)
		}
	}
}

func TestBuildArgs(t *testing.T) {
	cases := []struct {
		osArgs []string
		want []string
	}{
		{
			[]string{"Filename", "-b", "foo", "--build-arg", "bar", "--build-arg=foobar"},
			[]string{"foo", "bar", "foobar"},
		},
	}
	for _, c := range cases {
		got := ParseOsArgs(c.osArgs)
		if ! reflect.DeepEqual(got.BuildArgs, c.want) {
			t.Errorf("ParseOsArgs %q != %q", got.BuildArgs, c.want)
		}
	}
}
