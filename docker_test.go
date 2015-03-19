package dexec

import "testing"
import "./testutils"

func TestDockerPresence(t *testing.T) {
	cases := []struct {
		rawVersion string
		want bool
	}{
		{"Docker version 1.5.0, build abcdef0", true},
		{"Docker version x.y.z, build abcdef0", false},
		{"Mangled version string", false},
	}
	for _, c := range cases {
		defer testutils.Patch(&GetRawDockerVersion, func() string {
			return c.rawVersion
		}).Restore()

		got := IsDockerPresent()
		if got != c.want {
			t.Errorf("isDockerPresent() for version %q == %q, want %q", c.rawVersion, got, c.want)
		}
	}
}
