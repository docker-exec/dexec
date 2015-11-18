package cli

import (
	"reflect"
	"testing"
)

type OptionData struct {
	first  string
	second string
}

type WantedData struct {
	optionType   OptionType
	value        string
	chomped      int
	errorMessage string
}

func TestUnknown(t *testing.T) {
	cases := []struct {
		opt  OptionData
		want WantedData
	}{
		{
			OptionData{"-bad", "Unknown option: -bad"},
			WantedData{None, "", 0, "Unknown option: -bad"},
		},
	}
	for _, c := range cases {
		gotOptionType, gotOptionValue, gotChomped, gotError := ArgToOption(c.opt.first, c.opt.second)
		if gotOptionType != c.want.optionType {
			t.Errorf("ParseOsArgs %q != %q", gotOptionType, c.want.optionType)
		} else if gotOptionValue != c.want.value {
			t.Errorf("ParseOsArgs %q != %q", gotOptionValue, c.want.value)
		} else if gotChomped != c.want.chomped {
			t.Errorf("ParseOsArgs %q != %q", gotChomped, c.want.chomped)
		} else if gotError.Error() != c.want.errorMessage {
			t.Errorf("ParseOsArgs %q != %q", gotError.Error(), c.want.errorMessage)
		}
	}
}

func TestGet(t *testing.T) {
	cases := []struct {
		opt  OptionData
		want WantedData
	}{
		{
			OptionData{"foo.bar", ""},
			WantedData{Source, "foo.bar", 1, ""},
		},
		{
			OptionData{"-b", "foo"},
			WantedData{BuildArg, "foo", 2, ""},
		},
		{
			OptionData{"--build-arg", "foo"},
			WantedData{BuildArg, "foo", 2, ""},
		},
		{
			OptionData{"--build-arg=foo", ""},
			WantedData{BuildArg, "foo", 1, ""},
		},
		{
			OptionData{"-a", "foo"},
			WantedData{Arg, "foo", 2, ""},
		},
		{
			OptionData{"--arg", "foo"},
			WantedData{Arg, "foo", 2, ""},
		},
		{
			OptionData{"--arg=foo", ""},
			WantedData{Arg, "foo", 1, ""},
		},
		{
			OptionData{"-i", "foo"},
			WantedData{Include, "foo", 2, ""},
		},
		{
			OptionData{"--include", "foo"},
			WantedData{Include, "foo", 2, ""},
		},
		{
			OptionData{"--include=foo", ""},
			WantedData{Include, "foo", 1, ""},
		},
		{
			OptionData{"-m", "foo"},
			WantedData{SpecifyImage, "foo", 2, ""},
		},
		{
			OptionData{"--image", "foo"},
			WantedData{SpecifyImage, "foo", 2, ""},
		},
		{
			OptionData{"--image=foo", ""},
			WantedData{SpecifyImage, "foo", 1, ""},
		},
		{
			OptionData{"--help", ""},
			WantedData{HelpFlag, "", 1, ""},
		},
		{
			OptionData{"-h", ""},
			WantedData{HelpFlag, "", 1, ""},
		},
		{
			OptionData{"--version", ""},
			WantedData{VersionFlag, "", 1, ""},
		},
		{
			OptionData{"-v", ""},
			WantedData{VersionFlag, "", 1, ""},
		},
	}
	for _, c := range cases {
		gotOptionType, gotOptionValue, gotChomped, _ := ArgToOption(c.opt.first, c.opt.second)
		if gotOptionType != c.want.optionType {
			t.Errorf("ParseOsArgs %d != %d", gotOptionType, c.want.optionType)
		} else if gotOptionValue != c.want.value {
			t.Errorf("ParseOsArgs %s != %s", gotOptionValue, c.want.value)
		} else if gotChomped != c.want.chomped {
			t.Errorf("ParseOsArgs %d != %d", gotChomped, c.want.chomped)
		}
	}
}

func TestFilename(t *testing.T) {
	cases := []struct {
		osArgs []string
		want   string
	}{
		{[]string{"foo"}, "foo"},
		{[]string{"bar"}, "bar"},
	}
	for _, c := range cases {
		got := ParseOsArgs(c.osArgs)
		if got.Filename != c.want {
			t.Errorf("ParseOsArgs %q != %q", got.Filename, c.want)
		}
	}
}

func TestTargetDir(t *testing.T) {
	cases := []struct {
		osArgs []string
		want   []string
	}{
		{[]string{"filename", "-C", "foo"}, []string{"foo"}},
	}
	for _, c := range cases {
		got := ParseOsArgs(c.osArgs)
		if !reflect.DeepEqual(got.Options[TargetDir], c.want) {
			t.Errorf("ParseOsArgs %q != %q", got.Options[TargetDir], c.want)
		}
	}
}

func TestSources(t *testing.T) {
	cases := []struct {
		osArgs []string
		want   []string
	}{
		{
			[]string{"filename", "foo.cpp", "bar.java", "foo.bar.scala", "bar-foo.groovy"},
			[]string{"foo.cpp", "bar.java", "foo.bar.scala", "bar-foo.groovy"},
		},
	}
	for _, c := range cases {
		got := ParseOsArgs(c.osArgs)
		if !reflect.DeepEqual(got.Options[Source], c.want) {
			t.Errorf("ParseOsArgs %q != %q", got.Options[Source], c.want)
		}
	}
}

func TestSpecifyImage(t *testing.T) {
	cases := []struct {
		osArgs []string
		want   []string
	}{
		{
			[]string{"filename", "-m", "foo", "--image", "bar", "--image=foobar"},
			[]string{"foo", "bar", "foobar"},
		},
	}
	for _, c := range cases {
		got := ParseOsArgs(c.osArgs)
		if !reflect.DeepEqual(got.Options[SpecifyImage], c.want) {
			t.Errorf("ParseOsArgs %q != %q", got.Options[SpecifyImage], c.want)
		}
	}
}

func TestSpecifyExtension(t *testing.T) {
	cases := []struct {
		osArgs []string
		want   []string
	}{
		{
			[]string{"filename", "-e", "foo", "--extension", "bar", "--extension=foobar"},
			[]string{"foo", "bar", "foobar"},
		},
	}
	for _, c := range cases {
		got := ParseOsArgs(c.osArgs)
		if !reflect.DeepEqual(got.Options[Extension], c.want) {
			t.Errorf("ParseOsArgs %q != %q", got.Options[Extension], c.want)
		}
	}
}

func TestIncludes(t *testing.T) {
	cases := []struct {
		osArgs []string
		want   []string
	}{
		{
			[]string{"filename", "-i", "foo", "--include", "bar", "--include=foobar"},
			[]string{"foo", "bar", "foobar"},
		},
	}
	for _, c := range cases {
		got := ParseOsArgs(c.osArgs)
		if !reflect.DeepEqual(got.Options[Include], c.want) {
			t.Errorf("ParseOsArgs %q != %q", got.Options[Include], c.want)
		}
	}
}

func TestArgs(t *testing.T) {
	cases := []struct {
		osArgs []string
		want   []string
	}{
		{
			[]string{"filename", "-a", "foo", "--arg", "bar", "--arg=foobar"},
			[]string{"foo", "bar", "foobar"},
		},
	}
	for _, c := range cases {
		got := ParseOsArgs(c.osArgs)
		if !reflect.DeepEqual(got.Options[Arg], c.want) {
			t.Errorf("ParseOsArgs %q != %q", got.Options[Arg], c.want)
		}
	}
}

func TestBuildArgs(t *testing.T) {
	cases := []struct {
		osArgs []string
		want   []string
	}{
		{
			[]string{"filename", "-b", "foo", "--build-arg", "bar", "--build-arg=foobar"},
			[]string{"foo", "bar", "foobar"},
		},
	}
	for _, c := range cases {
		got := ParseOsArgs(c.osArgs)
		if !reflect.DeepEqual(got.Options[BuildArg], c.want) {
			t.Errorf("ParseOsArgs %q != %q", got.Options[BuildArg], c.want)
		}
	}
}

func TestOrdering(t *testing.T) {
	cases := []struct {
		osArgs []string
		want   map[OptionType][]string
	}{
		{
			[]string{
				"filename",
				"-C", "~/foo",
				"source1.foo",
				"-b", "b_foo",
				"-i", "i_foo",
				"source2.foo",
				"--arg=a_foobar",
				"source3.foo",
				"--build-arg", "b_bar",
				"--include", "i_bar",
				"source4.foo",
				"--arg", "a_bar",
				"source5.foo",
				"--build-arg=b_foobar",
				"--include=i_foobar",
				"source6.foo",
				"-a", "a_foo",
			},
			map[OptionType][]string{
				Arg: {
					"a_foobar", "a_bar", "a_foo",
				},
				BuildArg: {
					"b_foo", "b_bar", "b_foobar",
				},
				Include: {
					"i_foo", "i_bar", "i_foobar",
				},
				Source: {
					"source1.foo",
					"source2.foo",
					"source3.foo",
					"source4.foo",
					"source5.foo",
					"source6.foo",
				},
				TargetDir: {
					"~/foo",
				},
			},
		},
	}
	for _, c := range cases {
		got := ParseOsArgs(c.osArgs)
		if !reflect.DeepEqual(got.Options[BuildArg], c.want[BuildArg]) {
			t.Errorf("ParseOsArgs %q != %q", got.Options[BuildArg], c.want[BuildArg])
		} else if !reflect.DeepEqual(got.Options[Arg], c.want[Arg]) {
			t.Errorf("ParseOsArgs %q != %q", got.Options[Arg], c.want[Arg])
		} else if !reflect.DeepEqual(got.Options[Source], c.want[Source]) {
			t.Errorf("ParseOsArgs %q != %q", got.Options[Source], c.want[Source])
		}
	}
}
