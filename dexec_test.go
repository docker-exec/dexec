package main

import (
	"reflect"
	"testing"
)

func TestBuildVolumeArgs(t *testing.T) {
	cases := []struct {
		path        string
		targets     []string
		wantVolumes []string
	}{
		{"/foo", []string{"bar"}, []string{"-v", "/foo/bar:/tmp/dexec/build/bar"}},
	}
	for _, c := range cases {
		gotVolumes := BuildVolumeArgs(c.path, c.targets)
		if !reflect.DeepEqual(gotVolumes, c.wantVolumes) {
			t.Errorf("BuildVolumeArgs(%q, %q) %q != %q", c.path, c.targets, gotVolumes, c.wantVolumes)
		}
	}
}

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

func TestLookupImageByOverride(t *testing.T) {
	cases := []struct {
		image         string
		extension     string
		wantExtension string
		wantImage     string
		wantVersion   string
	}{
		{"dexec/cpp", "c", "c", "cpp", "latest"},
		{"dexec/some-language", "ext", "ext", "some-language", "latest"},
		{"dexec/some-language:1.2.3", "ext", "ext", "some-language", "1.2.3"},
	}
	for _, c := range cases {
		got := LookupImageByOverride(c.image, c.extension)
		if got.extension != c.wantExtension {
			t.Errorf("LookupImageByOverride(%q, %q) %q != %q", c.image, c.extension, got.image, c.wantImage)
		} else if got.extension != c.wantExtension {
			t.Errorf("LookupImageByOverride(%q, %q) %q != %q", c.image, c.extension, got.extension, c.wantExtension)
		} else if got.version != c.wantVersion {
			t.Errorf("LookupImageByOverride(%q, %q) %q != %q", c.image, c.extension, got.version, c.wantVersion)
		}
	}
}

func TestLookupImageByExtension(t *testing.T) {
	cases := []struct {
		extension     string
		wantExtension string
		wantImage     string
		wantVersion   string
	}{
		{"c", "c", "dexec/lang-c", "1.0.1"},
		{"clj", "clj", "dexec/lang-clojure", "1.0.0"},
		{"coffee", "coffee", "dexec/lang-coffee", "1.0.1"},
		{"cpp", "cpp", "dexec/lang-cpp", "1.0.1"},
		{"cs", "cs", "dexec/lang-csharp", "1.0.1"},
		{"d", "d", "dexec/lang-d", "1.0.0"},
		{"erl", "erl", "dexec/lang-erlang", "1.0.0"},
		{"fs", "fs", "dexec/lang-fsharp", "1.0.1"},
		{"go", "go", "dexec/lang-go", "1.0.0"},
		{"groovy", "groovy", "dexec/lang-groovy", "1.0.0"},
		{"hs", "hs", "dexec/lang-haskell", "1.0.0"},
		{"java", "java", "dexec/lang-java", "1.0.1"},
		{"lisp", "lisp", "dexec/lang-lisp", "1.0.0"},
		{"lua", "lua", "dexec/lang-lua", "1.0.0"},
		{"js", "js", "dexec/lang-node", "1.0.1"},
		{"m", "m", "dexec/lang-objc", "1.0.0"},
		{"ml", "ml", "dexec/lang-ocaml", "1.0.0"},
		{"nim", "nim", "dexec/lang-nim", "1.0.0"},
		{"p6", "p6", "dexec/lang-perl6", "1.0.0"},
		{"pl", "pl", "dexec/lang-perl", "1.0.1"},
		{"php", "php", "dexec/lang-php", "1.0.0"},
		{"py", "py", "dexec/lang-python", "1.0.1"},
		{"r", "r", "dexec/lang-r", "1.0.0"},
		{"rkt", "rkt", "dexec/lang-racket", "1.0.0"},
		{"rb", "rb", "dexec/lang-ruby", "1.0.0"},
		{"rs", "rs", "dexec/lang-rust", "1.0.0"},
		{"scala", "scala", "dexec/lang-scala", "1.0.0"},
		{"sh", "sh", "dexec/lang-bash", "1.0.0"},
	}
	for _, c := range cases {
		got := LookupImageByExtension(c.extension)
		if got.image != c.wantImage {
			t.Errorf("TestLookupExtensionByImage(%q) %q != %q", c.extension, got.image, c.wantImage)
		} else if got.extension != c.wantExtension {
			t.Errorf("TestLookupExtensionByImage(%q) %q != %q", c.extension, got.extension, c.wantExtension)
		} else if got.version != c.wantVersion {
			t.Errorf("TestLookupExtensionByImage(%q) %q != %q", c.extension, got.version, c.wantVersion)
		}
	}
}
