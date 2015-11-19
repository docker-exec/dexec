package util

import (
	"reflect"
	"testing"
)

func TestSanitisePath(t *testing.T) {
	cases := []struct {
		path     string
		platform string
		want     string
	}{
		{"/Users/foo/bar", "darin", "/Users/foo/bar"},
		{"/home/foo/bar", "linux", "/home/foo/bar"},
		{"C:\\Users\\foo\\bar", "windows", "/c/Users/foo/bar"},
	}
	for _, c := range cases {
		gotSanitisedPath := SanitisePath(c.path, c.platform)
		if gotSanitisedPath != c.want {
			t.Errorf("SanitisedPath %q != %q", gotSanitisedPath, c.want)
		}
	}
}

func TestAddPrefix(t *testing.T) {
	cases := []struct {
		inSlice []string
		prefix  string
		want    []string
	}{
		{
			[]string{"foo", "bar", "foobar"},
			"prefix",
			[]string{"prefix", "foo", "prefix", "bar", "prefix", "foobar"},
		},
	}
	for _, c := range cases {
		got := AddPrefix(c.inSlice, c.prefix)
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("AddPrefix(%q, %q) %q != %q", c.inSlice, c.prefix, got, c.want)
		}
	}
}

func TestJoinStringSlices(t *testing.T) {
	cases := []struct {
		inSlices [][]string
		want     []string
	}{
		{
			[][]string{{"foo"}, {"bar"}, {"foobar"}},
			[]string{"foo", "bar", "foobar"},
		},
	}
	for _, c := range cases {
		got := JoinStringSlices(c.inSlices...)
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("JoinStringSlices(%q) %q != %q", c.inSlices, got, c.want)
		}
	}
}

func TestExtractFileExtension(t *testing.T) {
	cases := []struct {
		filename  string
		extension string
	}{
		{"foo.bar", "bar"},
		{"foo.bar.foobar", "foobar"},
	}
	for _, c := range cases {
		gotExtension := ExtractFileExtension(c.filename)
		if gotExtension != c.extension {
			t.Errorf("ExtractFileExtension %q != %q", gotExtension, c.extension)
		}
	}
}
