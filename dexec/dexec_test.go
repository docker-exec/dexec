package dexec

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
		{"/foo", []string{"bar"}, []string{"/foo/bar:/tmp/dexec/build/bar"}},
	}
	for _, c := range cases {
		gotVolumes := BuildVolumeArgs(c.path, c.targets)
		if !reflect.DeepEqual(gotVolumes, c.wantVolumes) {
			t.Errorf("BuildVolumeArgs(%q, %q) %q != %q", c.path, c.targets, gotVolumes, c.wantVolumes)
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
		wantError     error
	}{
		{"dexec/cpp", "c", "c", "dexec/cpp", "latest", nil},
		{"dexec/some-language", "ext", "ext", "dexec/some-language", "latest", nil},
		{"dexec/some-language:1.2.3", "ext", "ext", "dexec/some-language", "1.2.3", nil},
	}
	for _, c := range cases {
		got, err := LookupImageByOverride(c.image, c.extension)
		if got.Image != c.wantImage {
			t.Errorf("LookupImageByOverride(%q, %q) %q != %q", c.image, c.extension, got.Image, c.wantImage)
		} else if got.Extension != c.wantExtension {
			t.Errorf("LookupImageByOverride(%q, %q) %q != %q", c.image, c.extension, got.Extension, c.wantExtension)
		} else if got.Version != c.wantVersion {
			t.Errorf("LookupImageByOverride(%q, %q) %q != %q", c.image, c.extension, got.Version, c.wantVersion)
		} else if err != c.wantError {
			t.Errorf("LookupImageByOverride(%q, %q) %q != %q", c.image, c.extension, err, c.wantError)
		}
	}
}

func TestLookupImageByExtension(t *testing.T) {
	cases := []struct {
		extension     string
		wantExtension string
		wantImage     string
		wantVersion   string
		wantError     error
	}{
		{"c", "c", "dexec/lang-c", "1.0.2", nil},
		{"clj", "clj", "dexec/lang-clojure", "1.0.1", nil},
		{"coffee", "coffee", "dexec/lang-coffee", "1.0.2", nil},
		{"cpp", "cpp", "dexec/lang-cpp", "1.0.2", nil},
		{"cs", "cs", "dexec/lang-csharp", "1.0.2", nil},
		{"d", "d", "dexec/lang-d", "1.0.1", nil},
		{"erl", "erl", "dexec/lang-erlang", "1.0.1", nil},
		{"fs", "fs", "dexec/lang-fsharp", "1.0.2", nil},
		{"go", "go", "dexec/lang-go", "1.0.1", nil},
		{"groovy", "groovy", "dexec/lang-groovy", "1.0.1", nil},
		{"hs", "hs", "dexec/lang-haskell", "1.0.1", nil},
		{"java", "java", "dexec/lang-java", "1.0.3", nil},
		{"lisp", "lisp", "dexec/lang-lisp", "1.0.1", nil},
		{"lua", "lua", "dexec/lang-lua", "1.0.1", nil},
		{"js", "js", "dexec/lang-node", "1.0.2", nil},
		{"m", "m", "dexec/lang-objc", "1.0.2", nil},
		{"ml", "ml", "dexec/lang-ocaml", "1.0.1", nil},
		{"nim", "nim", "dexec/lang-nim", "1.0.1", nil},
		{"p6", "p6", "dexec/lang-perl6", "1.0.1", nil},
		{"pl", "pl", "dexec/lang-perl", "1.0.2", nil},
		{"php", "php", "dexec/lang-php", "1.0.1", nil},
		{"py", "py", "dexec/lang-python", "1.0.2", nil},
		{"r", "r", "dexec/lang-r", "1.0.1", nil},
		{"rkt", "rkt", "dexec/lang-racket", "1.0.1", nil},
		{"rb", "rb", "dexec/lang-ruby", "1.0.2", nil},
		{"rs", "rs", "dexec/lang-rust", "1.0.1", nil},
		{"scala", "scala", "dexec/lang-scala", "1.0.1", nil},
		{"sh", "sh", "dexec/lang-bash", "1.0.1", nil},
	}
	for _, c := range cases {
		got, err := LookupImageByExtension(c.extension)
		if got.Image != c.wantImage {
			t.Errorf("TestLookupExtensionByImage(%q) %q != %q", c.extension, got.Image, c.wantImage)
		} else if got.Extension != c.wantExtension {
			t.Errorf("TestLookupExtensionByImage(%q) %q != %q", c.extension, got.Extension, c.wantExtension)
		} else if got.Version != c.wantVersion {
			t.Errorf("TestLookupExtensionByImage(%q) %q != %q", c.extension, got.Version, c.wantVersion)
		} else if err != c.wantError {
			t.Errorf("TestLookupExtensionByImage(%q) %q != %q", c.extension, err, c.wantError)
		}
	}
}
