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

func TestLookupExtensionByImage(t *testing.T) {
	cases := []struct {
		extension string
		image     string
	}{
		{"c", "c"},
		{"clj", "clojure"},
		{"coffee", "coffee"},
		{"cpp", "cpp"},
		{"cs", "csharp"},
		{"d", "d"},
		{"erl", "erlang"},
		{"fs", "fsharp"},
		{"go", "go"},
		{"groovy", "groovy"},
		{"hs", "haskell"},
		{"java", "java"},
		{"lisp", "lisp"},
		{"js", "node"},
		{"m", "objc"},
		{"ml", "ocaml"},
		{"pl", "perl"},
		{"php", "php"},
		{"py", "python"},
		{"rkt", "racket"},
		{"rb", "ruby"},
		{"rs", "rust"},
		{"scala", "scala"},
		{"sh", "bash"},
	}
	for _, c := range cases {
		gotImage := LookupExtensionByImage(c.extension)
		if gotImage != c.image {
			t.Errorf("TestLookupExtensionByImage %q != %q", gotImage, c.image)
		}
	}
}
