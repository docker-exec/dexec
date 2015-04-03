package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
)

const dexecPath = "/tmp/dexec/build"
const dexecImageTemplate = "dexec/%s"
const dexecVolumeTemplate = "%s:%s:ro"

func AddPrefix(inSlice []string, prefix string) []string {
	outSlice := []string{}
	for _, option := range inSlice {
		outSlice = append(outSlice, []string{prefix, option}...)
	}
	return outSlice
}

func JoinStringSlices(slices ...[]string) []string {
	var outSlice []string
	for _, slice := range slices {
		outSlice = append(outSlice, slice...)
	}
	return outSlice
}


var DockerVersion = func() string {
	out, err := exec.Command("docker", "-v").Output()
	if err != nil {
		panic(err.Error())
	} else {
		return string(out)
	}
}

var DockerInfo = func() string {
	out, err := exec.Command("docker", "info").Output()
	if err != nil {
		panic(err.Error())
	} else {
		return string(out)
	}
}

func ExtractDockerVersion(version string) (int, int, int) {
	dockerVersionPattern := regexp.MustCompile(`^Docker version (\d+)\.(\d+)\.(\d+), build [a-z0-9]{7}`)

	if dockerVersionPattern.MatchString(version) {
		match := dockerVersionPattern.FindStringSubmatch(version)
		major, _ := strconv.Atoi(match[1])
		minor, _ := strconv.Atoi(match[2])
		patch, _ := strconv.Atoi(match[3])
		return major, minor, patch
	} else {
		panic("Did not match Docker version string")
	}
}

func IsDockerPresent() bool {
	present := true
	defer func() {
		if r := recover(); r != nil {
			present = false
		}
	}()
	ExtractDockerVersion(DockerVersion())
	return present
}

func IsDockerRunning() bool {
	running := true
	defer func() {
		if r := recover(); r != nil {
			running = false
		}
	}()
	DockerInfo()
	return running
}


func RunAnonymousContainer(image string, extraDockerArgs []string, entrypointArgs []string) {
	baseDockerArgs := []string{"run", "--rm"}
	imageDockerArgs := []string{"-t", image}

	out := exec.Command(
		"docker",
		JoinStringSlices(
			baseDockerArgs,
			extraDockerArgs,
			imageDockerArgs,
			entrypointArgs,
		)...,
	)
	out.Stdin = os.Stdin
	out.Stdout = os.Stdout
	out.Stderr = os.Stderr
	out.Run()
}

func RunDexecContainer(image string, options map[OptionType][]string) {
	absPath, _ := filepath.Abs(".")

	dockerArgs := []string{
		"-v",
		fmt.Sprintf(dexecVolumeTemplate, absPath, dexecPath),
	}

	entrypointArgs := JoinStringSlices(
		options[Source],
		AddPrefix(options[BuildArg], "-b"),
		AddPrefix(options[Arg], "-a"),
	)

	RunAnonymousContainer(
		fmt.Sprintf(dexecImageTemplate, image),
		dockerArgs,
		entrypointArgs,
	)
}
