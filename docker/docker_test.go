package docker

import (
	"reflect"
	"testing"
)

// Restorer holds a function that can be used
// to restore some previous state.
type Restorer func()

// Restore restores some previous state.
func (r Restorer) Restore() {
	r()
}

// Patch sets the value pointed to by the given destination to the given
// value, and returns a function to restore it to its original value.  The
// value must be assignable to the element type of the destination.
func Patch(dest, value interface{}) Restorer {
	destv := reflect.ValueOf(dest).Elem()
	oldv := reflect.New(destv.Type()).Elem()
	oldv.Set(destv)
	valuev := reflect.ValueOf(value)
	if !valuev.IsValid() {
		valuev = reflect.Zero(destv.Type())
	}
	destv.Set(valuev)
	return func() {
		destv.Set(oldv)
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

func TestExtractDockerVersion(t *testing.T) {
	cases := []struct {
		version string
		want    [3]int
	}{
		{"Docker version 1.5.0, build abcdef0", [3]int{1, 5, 0}},
	}
	for _, c := range cases {
		major, minor, patch := ExtractDockerVersion(c.version)
		if major != c.want[0] || minor != c.want[1] || patch != c.want[2] {
			t.Errorf("ExtractDockerVersion(%q) %q.%q.%q != %q.%q.%q", c.version, major, minor, patch, c.want[0], c.want[1], c.want[2])
		}
	}
}

func TestIsDockerPresent(t *testing.T) {
	cases := []struct {
		version string
		want    bool
	}{
		{"Docker version 1.5.0, build abcdef0", true},
		{"Docker version x.y.z, build abcdef0", false},
		{"Mangled version string", false},
	}
	for _, c := range cases {
		defer Patch(&DockerVersion, func() string {
			return c.version
		}).Restore()

		got := IsDockerPresent()
		if got != c.want {
			t.Errorf("IsDockerPresent() for version %q == %v, want %v", c.version, got, c.want)
		}
	}
}

func TestIsDockerRunning(t *testing.T) {
	cases := []struct {
		output string
		want   bool
	}{
		{"Docker info string", true},
	}
	for _, c := range cases {
		defer Patch(&DockerInfo, func() string {
			return c.output
		}).Restore()

		got := IsDockerRunning()
		if got != c.want {
			t.Errorf("IsDockerRunning() for info string %q == %v, want %v", c.output, got, c.want)
		}
	}
}
