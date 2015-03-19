package dexec

import "testing"

func TestDockerPresence(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"Docker version 1.5.0, build a8a31ef", (true, "1.5.0")},
	}
	for _, c := range cases {
		present, version := isDockerPresent(c.in)
		if got != c.want {
			t.Errorf("Reverse(%q) == %q, want %q", c.in, got, c.want)
		}
	}
}