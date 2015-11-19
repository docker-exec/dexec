package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"

	"github.com/docker-exec/dexec/cli"
	"github.com/docker-exec/dexec/dexec"
	"github.com/docker-exec/dexec/util"
	"github.com/fsouza/go-dockerclient"
)

// RunDexecContainer runs an anonymous Docker container with a Docker Exec
// image, mounting the specified sources and includes and passing the
// list of sources and arguments to the entrypoint.
func RunDexecContainer(cliParser cli.CLI) {
	options := cliParser.Options
	dexecImage, err := dexec.ImageFromOptions(options)

	if err != nil {
		log.Fatal(err)
	}

	dockerImage := fmt.Sprintf("%s:%s", dexecImage.Image, dexecImage.Version)
	updateImage := len(options[cli.UpdateFlag]) > 0

	client, err := docker.NewClientFromEnv()

	if err != nil {
		log.Fatal(err)
	}

	if err := dexec.FetchImage(
		dexecImage.Image,
		dexecImage.Version,
		updateImage,
		client); err != nil {
		log.Fatal(err)
	}

	if useStdin := len(options[cli.Source]) == 0; useStdin {
		lines := util.ReadStdin("Enter your code. Ctrl-D to exit", "<Ctrl-D>")
		tmpFile := fmt.Sprintf("%s.%s", uuid.NewV4(), dexecImage.Extension)

		util.WriteFile(tmpFile, []byte(strings.Join(lines, "\n")))
		defer func() {
			util.DeleteFile(tmpFile)
		}()

		options[cli.Source] = []string{tmpFile}
	}

	var sourceBasenames []string
	for _, source := range options[cli.Source] {
		basename, _ := dexec.ExtractBasenameAndPermission(source)
		sourceBasenames = append(sourceBasenames, []string{basename}...)
	}

	entrypointArgs := util.JoinStringSlices(
		sourceBasenames,
		util.AddPrefix(options[cli.BuildArg], "-b"),
		util.AddPrefix(options[cli.Arg], "-a"),
	)

	container, err := client.CreateContainer(docker.CreateContainerOptions{
		Config: &docker.Config{
			Image: dockerImage,
			Cmd:   entrypointArgs,
		},
	})

	if err != nil {
		log.Fatal(err)
	}

	if err = client.StartContainer(container.ID, &docker.HostConfig{
		Binds: dexec.BuildVolumeArgs(
			util.RetrievePath(options[cli.TargetDir]),
			append(options[cli.Source], options[cli.Include]...)),
	}); err != nil {
		log.Fatal(err)
	}

	if err = client.AttachToContainer(docker.AttachToContainerOptions{
		Container:    container.ID,
		OutputStream: os.Stdout,
		ErrorStream:  os.Stderr,
		Stream:       true,
		Stdout:       true,
		Stderr:       true,
		Logs:         true,
	}); err != nil {
		log.Fatal(err)
	}

	if err = client.RemoveContainer(docker.RemoveContainerOptions{
		ID: container.ID,
	}); err != nil {
		log.Fatal(err)
	}
}

func validate(cliParser cli.CLI) bool {
	options := cliParser.Options

	hasVersionFlag := len(options[cli.VersionFlag]) == 1
	hasExtension := len(options[cli.Extension]) == 1
	hasImage := len(options[cli.Image]) == 1
	hasSources := len(options[cli.Source]) > 0

	if hasSources || hasImage || hasExtension {
		return true
	}

	if hasVersionFlag {
		cli.DisplayVersion(cliParser.Filename)
		return false
	}

	cli.DisplayHelp(cliParser.Filename)
	return false
}

func validateDocker() error {
	client, err := docker.NewClientFromEnv()
	if err != nil {
		return err
	}

	ping := make(chan error, 1)
	go func() {
		ping <- client.Ping()
	}()

	select {
	case err := <-ping:
		return err
	case <-time.After(5 * time.Second):
		return fmt.Errorf("Request to Docker host timed out")
	}
}

func main() {
	cliParser := cli.ParseOsArgs(os.Args)

	if validate(cliParser) {
		if err := validateDocker(); err != nil {
			log.Fatal(err)
		} else {
			RunDexecContainer(cliParser)
		}
	}
}
