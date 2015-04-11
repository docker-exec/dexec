package main

import "testing"

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
		{"c", "c", "c", "1.0.0"},
		{"clj", "clj", "clojure", "1.0.0"},
		{"coffee", "coffee", "coffee", "1.0.0"},
		{"cpp", "cpp", "cpp", "1.0.0"},
		{"cs", "cs", "csharp", "1.0.0"},
		{"d", "d", "d", "1.0.0"},
		{"erl", "erl", "erlang", "1.0.0"},
		{"fs", "fs", "fsharp", "1.0.0"},
		{"go", "go", "go", "1.0.0"},
		{"groovy", "groovy", "groovy", "1.0.0"},
		{"hs", "hs", "haskell", "1.0.0"},
		{"java", "java", "java", "1.0.0"},
		{"lisp", "lisp", "lisp", "1.0.0"},
		{"lua", "lua", "lua", "latest"},
		{"js", "js", "node", "1.0.0"},
		{"m", "m", "objc", "1.0.0"},
		{"ml", "ml", "ocaml", "1.0.0"},
		{"nim", "nim", "nim", "latest"},
		{"p6", "p6", "perl6", "1.0.0"},
		{"pl", "pl", "perl", "1.0.1"},
		{"php", "php", "php", "1.0.0"},
		{"py", "py", "python", "1.0.0"},
		{"rkt", "rkt", "racket", "1.0.0"},
		{"rb", "rb", "ruby", "1.0.0"},
		{"rs", "rs", "rust", "1.0.0"},
		{"scala", "scala", "scala", "1.0.0"},
		{"sh", "sh", "bash", "1.0.0"},
	}
	for _, c := range cases {
		got := LookupImageByExtension(c.extension)
		if got.extension != c.wantExtension {
			t.Errorf("TestLookupExtensionByImage(%q) %q != %q", c.extension, got.image, c.wantImage)
		} else if got.extension != c.wantExtension {
			t.Errorf("TestLookupExtensionByImage(%q) %q != %q", c.extension, got.extension, c.wantExtension)
		} else if got.version != c.wantVersion {
			t.Errorf("TestLookupExtensionByImage(%q) %q != %q", c.extension, got.version, c.wantVersion)
		}
	}
}
