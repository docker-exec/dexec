package dexec

import (
	"code.google.com/p/go-uuid/uuid"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
)

var GetRawDockerVersion = func() string {
	out, err := exec.Command("docker", "-v").Output()
	if err != nil {
		panic(err.Error())
	} else {
		return string(out)
	}
}

func ExtractDockerVersion(rawVersion string) (int, int, int) {
	dockerVersionPattern := regexp.MustCompile(`^Docker version (\d+)\.(\d+)\.(\d+), build [a-z0-9]{7}`)

	if dockerVersionPattern.MatchString(rawVersion) {
		match := dockerVersionPattern.FindStringSubmatch(rawVersion)
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

	ExtractDockerVersion(GetRawDockerVersion())

	return present
}

func IsDockerRunning() bool {
	out := exec.Command("docker", "info")
	return out.Run() != nil
}

func RunAnonymousContainer(args ...string) {
	newArgs := append([]string{"run", "-t", "--rm"}, args...)
	out := exec.Command("docker", newArgs...)
	out.Stdin = os.Stdin
	out.Stdout = os.Stdout
	out.Stderr = os.Stderr
	out.Run()
}

func RunDexecContainer(language string, sourcefile string, entrypointargs ...string) {
	workdir := fmt.Sprintf("/tmp/%s", uuid.New())
	abssourcefile, _ := filepath.Abs(sourcefile)

	RunAnonymousContainer(
		append(
			[]string{
				"-w", workdir,
				"-v", fmt.Sprintf("%s:%s/%s", abssourcefile, workdir, sourcefile),
				fmt.Sprintf("dexec/%s", language),
				fmt.Sprintf("%s/%s", workdir, sourcefile)},
			entrypointargs...,
		)...,
	)
}