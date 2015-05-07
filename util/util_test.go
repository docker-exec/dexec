package util

import (
	"reflect"
	"testing"
)

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
