package main

import (
	"testing"
	"reflect"
)


func TestFilenameR(t *testing.T) {
	cases := []struct {
		osArgs []string
		want string
	}{
		{ []string{"foo"}, "foo" },
		{ []string{"bar"}, "bar" },
	}
	for _, c := range cases {
		got := ParseOsArgsR(c.osArgs)
		if got.FileName != c.want {
			t.Errorf("ParseOsArgsR %q != %q", got.FileName, c.want)
		}
	}
}

func TestArgsR(t *testing.T) {
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
		got := ParseOsArgsR(c.osArgs)
		if ! reflect.DeepEqual(got.Args, c.want) {
			t.Errorf("ParseOsArgs %q != %q", got.Args, c.want)
		}
	}
}
//
// func TestBuildArgs(t *testing.T) {
// 	cases := []struct {
// 		osArgs []string
// 		want []string
// 	}{
// 		{
// 			[]string{"Filename", "-b", "foo", "--build-arg", "bar", "--build-arg=foobar"},
// 			[]string{"foo", "bar", "foobar"},
// 		},
// 	}
// 	for _, c := range cases {
// 		got := ParseOsArgsR(c.osArgs)
// 		if ! reflect.DeepEqual(got.BuildArgs, c.want) {
// 			t.Errorf("ParseOsArgs %q != %q", got.BuildArgs, c.want)
// 		}
// 	}
// }
//
// func TestSources(t *testing.T) {
// 	cases := []struct {
// 		osArgs []string
// 		want []string
// 	}{
// 		{
// 			[]string{"Filename", "foo.cpp", "bar.java", "foo.bar.scala", "bar-foo.groovy"},
// 			[]string{"foo.cpp", "bar.java", "foo.bar.scala", "bar-foo.groovy"},
// 		},
// 	}
// 	for _, c := range cases {
// 		got := ParseOsArgsR(c.osArgs)
// 		if ! reflect.DeepEqual(got.Sources, c.want) {
// 			t.Errorf("ParseOsArgs %q != %q", got.Sources, c.want)
// 		}
// 	}
// }
